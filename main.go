package main

import (
	"chatApp/api"
	"chatApp/db"
	"chatApp/server"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "modernc.org/sqlite"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var currentUserId int = 1

func main() {
	fmt.Println("Hello, World")
	db.OpenConnectionToDatabase()
	db.IsDatabaseInitialized = true
	setupURLRoutes()
	log.Fatal(http.ListenAndServe(":80", nil))
}

func setupURLRoutes() {
	fileServer := http.FileServer(http.Dir("./content"))
	http.Handle("/content/", http.StripPrefix("/content/", fileServer))
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", websocketEndpoint)
	http.HandleFunc("/api/connectedUsers", api.GetConnectedUsers)
	http.HandleFunc("/admin", adminDashboard)
	http.HandleFunc("/api/verifyUserLoginData/", api.VerifyUserLoginData)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if server.AreLoginDataParamsSet(r) {
		http.SetCookie(w, &http.Cookie{Name: "logindata", Value: fmt.Sprintf("%v:%v", r.PostForm.Get("username"), r.PostForm.Get("password"))})
	}
	http.ServeFile(w, r, "./content/index.html")
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	session := server.Session{Connection: ws, Id: currentUserId, Name: fmt.Sprintf("User %v", currentUserId)}
	currentUserId++
	checkError(err)
	log.Printf("%v connected", session.Name)
	handle(session)
}

func handle(session server.Session) {
	server.InitializeSession(&session)
	server.SendJoinMessageForSession(&session)
	server.ConnectedUsers = append(server.ConnectedUsers, session)
	server.HandleSession(session)
}

func adminDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./content/admin.html")
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
