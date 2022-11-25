package server

import (
	"chatApp/db"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Session struct {
	Connection *websocket.Conn
	Id         int
	Name       string
}

type Message struct {
	Content  string `json:"content"`
	Sender   string `json:"sender"`
	SenderId int    `json:"sender_id"`
	Type     int    `json:"type"`
}

type ErrNoSessionFound struct {
	Where string
}

func (E ErrNoSessionFound) Error() string {
	return fmt.Sprintf("Found no session in %v", E.Where)
}

const (
	PINGMESSAGE   = 9
	SYSTEMMESSAGE = iota
	CHATMESSAGE
)

const (
	NOSESSIONFOUND = -1
)

var ConnectedUsers []Session = make([]Session, 0)

func InitializeSession(session *Session) {
	_, p, err := session.Connection.ReadMessage()
	checkError(err)
	message := string(p)
	username, password, _ := strings.Cut(message, ":")
	hashedPassword := GenerateHashForPassword(password)
	data := db.GetUserDataFromDatabaseBy("username", username)
	validatePasswordAndUsername(data, username, hashedPassword, session)
}

func GenerateHashForPassword(password string) uint32 {
	hashFunction := fnv.New32a()
	hashFunction.Write([]byte(password))
	return hashFunction.Sum32()
}

func validatePasswordAndUsername(data *db.User, username string, password uint32, session *Session) {
	if data.Username == username && data.Password == password {
		session.Id = data.Id
		db.UpdateValueOfUser("connected", true, data.Id)
		sendSuccessMessages(username, session)
	} else {
		msg := Message{Content: "NOTLOGGEDIN", Sender: "System", SenderId: 0, Type: SYSTEMMESSAGE}
		session.Connection.WriteJSON(msg)
	}
}

func sendSuccessMessages(username string, session *Session) {
	msg := Message{Content: "LOGGEDIN", Sender: "System", SenderId: 0, Type: SYSTEMMESSAGE}
	session.Connection.WriteJSON(msg)
	session.Name = username
	msg.Content = "Logged in succesfully"
	msg.Type = CHATMESSAGE
	session.Connection.WriteJSON(msg)
	msg.Content = fmt.Sprintf("Hello %v", session.Name)
	session.Connection.WriteJSON(msg)
}

func HandleSession(s Session) {
	keepConnectionAlive(s.Connection, 500*time.Millisecond)
	s.Connection.SetCloseHandler(createCloseHandlerFor(s))
	for {
		msg, err := listenForMessages(s)
		if err != nil {
			return
		}
		log.Printf("Message by %v: %v", s.Name, msg.Content)
		forwardMesage(msg)
	}
}

func keepConnectionAlive(c *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		for {
			if err := sendPingMessageTo(c); err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Since(lastResponse) > timeout {
				session := getSessionOfConnection(ConnectedUsers, c)
				SendLeftMessageForSession(&session)
				disconnectToClientConnection(c)
				return
			}
		}
	}()
}

func createCloseHandlerFor(s Session) func(code int, text string) error {
	return func(code int, text string) error {
		if i := getIndexOfSession(ConnectedUsers, s); i != NOSESSIONFOUND && (code == websocket.CloseGoingAway || code == websocket.CloseNoStatusReceived) {
			SendLeftMessageForSession(&s)
			ConnectedUsers = updateConnectedUsers(ConnectedUsers, i)
			db.UpdateValueOfUser("connected", false, s.Id)
		}
		return nil
	}
}

func listenForMessages(s Session) (msg Message, err error) {
	_, p, err := s.Connection.ReadMessage()
	msg = Message{Content: string(p), Sender: s.Name, SenderId: s.Id, Type: CHATMESSAGE}
	if err != nil {
		log.Printf("%v disconnected: %v", s.Name, err)
		ConnectedUsers = updateConnectedUsers(ConnectedUsers, getIndexOfSession(ConnectedUsers, s))
		return msg, err
	}
	return msg, nil
}

func forwardMesage(msg Message) {
	for index, c := range ConnectedUsers {
		deleteConnectionIfClosed(msg, index, c)
	}
}

func SendJoinMessageForSession(session *Session) {
	forwardMesage(Message{Content: fmt.Sprintf("%v joined", session.Name),
		Type: CHATMESSAGE, Sender: "System", SenderId: 0})
}

func SendLeftMessageForSession(session *Session) {
	leftMsg := Message{Content: fmt.Sprintf("%v left", session.Name),
		Sender: "System", SenderId: 0, Type: CHATMESSAGE}
	forwardMesage(leftMsg)
}

func updateConnectedUsers(connectedUsers []Session, index int) []Session {
	newArray := make([]Session, 0)
	for i := 0; i < len(connectedUsers); i++ {
		if index != i {
			newArray = append(newArray, connectedUsers[i])
		}
	}
	return newArray
}

func getIndexOfSession(connectedUsers []Session, s Session) int {
	for index, session := range connectedUsers {
		if s.Id == session.Id {
			return index
		}
	}
	return NOSESSIONFOUND
}

func deleteConnectionIfClosed(msgToSend Message, index int, c Session) {
	if err := c.Connection.WriteJSON(msgToSend); err == websocket.ErrCloseSent {
		log.Println(err)
		ConnectedUsers = updateConnectedUsers(ConnectedUsers, index)
		db.UpdateValueOfUser("connected", false, c.Id)
	}
}

func getSessionOfConnection(connectedUsers []Session, c *websocket.Conn) Session {
	for _, s := range connectedUsers {
		if s.Connection == c {
			return s
		}
	}
	return Session{Id: -1}
}

func sendPingMessageTo(c *websocket.Conn) error {
	session := getSessionOfConnection(ConnectedUsers, c)
	if session.Id == NOSESSIONFOUND {
		return ErrNoSessionFound{Where: "ConnectedUsers"}
	}
	log.Printf("pinging %v", session.Name)
	err := c.WriteMessage(PINGMESSAGE, []byte("keepalive"))
	if err != nil {
		return err
	}
	return nil
}

func disconnectToClientConnection(c *websocket.Conn) {
	log.Printf("Ping don't get response, disconnecting to %s", c.RemoteAddr())
	db.UpdateValueOfUser("connected", false, getSessionOfConnection(ConnectedUsers, c).Id)
	ConnectedUsers = updateConnectedUsers(ConnectedUsers, getIndexOfSession(ConnectedUsers, getSessionOfConnection(ConnectedUsers, c)))
	_ = c.Close()
}

func AreLoginDataParamsSet(r *http.Request) bool {
	return r.PostForm.Get("username") != "" && r.PostForm.Get("password") != ""
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
