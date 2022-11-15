package api

import (
	"chatApp/db"
	"encoding/json"
	"net/http"
)

func GetConnectedUsers(w http.ResponseWriter, r *http.Request) {
	users := db.GetConnectedUsers()
	jsonBytes, _ := json.Marshal(users)
	w.Write(jsonBytes)
}
