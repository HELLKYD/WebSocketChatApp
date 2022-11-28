package api

import (
	"chatApp/db"
	"chatApp/server"
	"encoding/json"
	"net/http"
	"strings"
)

func GetConnectedUsers(w http.ResponseWriter, r *http.Request) {
	users := db.GetConnectedUsers()
	jsonBytes, _ := json.Marshal(users)
	w.Write(jsonBytes)
}

type response struct {
	DataValid bool `json:"isValid"`
}

func VerifyUserLoginData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userData := strings.Split(r.Form.Get("loginData"), ":")
	username := userData[0]
	password := userData[1]
	hashedPassword := server.GenerateHashForPassword(password)
	dbUserData := db.GetUserDataFromDatabaseBy("username", "Administrator")
	if dbUserData.Password == hashedPassword && dbUserData.Username == username {
		response, _ := json.Marshal(response{true})
		w.Write(response)
	} else {
		response, _ := json.Marshal(response{false})
		w.Write(response)
	}
}
