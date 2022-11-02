package main

import (
	"chatApp/db"
	"chatApp/server"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
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
}

func homePage(w http.ResponseWriter, r *http.Request) {
	user := db.GetUserDataFromDatabaseBy("id", 0)
	fmt.Printf("user: %v\n", user)
	http.ServeFile(w, r, "./content/index.html")
}

func websocketEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	session := server.Session{Connection: ws, Id: currentUserId, Name: fmt.Sprintf("User %v", currentUserId)}
	currentUserId++
	checkError(err)
	log.Printf("%v connected", session.Name)
	server.InitializeSession(ws, &session)
	server.ConnectedUsers = append(server.ConnectedUsers, session)
	go server.HandleSession(session)
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
