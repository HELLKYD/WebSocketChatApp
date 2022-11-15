package db

import (
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	Connected bool   `json:"connected,omitempty"`
}

type validTypes interface {
	~string | ~int | ~bool
}

var DB *sql.DB

var IsDatabaseInitialized = false

func OpenConnectionToDatabase() {
	temp_db, err := sql.Open("sqlite", "./test.db")
	if err != nil {
		panic(err)
	}
	DB = temp_db
}

func GetUserDataFromDatabaseBy[T validTypes](columnName string, value T) *User {
	if IsDatabaseInitialized {
		return retrieveData(columnName, value)
	} else {
		OpenConnectionToDatabase()
		IsDatabaseInitialized = true
		return retrieveData(columnName, value)
	}
}

func retrieveData[T validTypes](columnName string, value T) *User {
	user := User{}
	data := DB.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE %v='%v';", columnName, value))
	data.Scan(&user.Id, &user.Username, &user.Password, &user.Connected)
	return &user
}

func UpdateValueOfUser[T validTypes](column string, newValue T, userId int) {
	DB.Exec(fmt.Sprintf("UPDATE users SET %v = %v WHERE id = %v;", column, newValue, userId))
}

func GetConnectedUsers() []User {
	users := make([]User, 0)
	data, err := DB.Query("SELECT id, username FROM users WHERE connected = true;")
	if err != nil {
		log.Println("Error while retrieving data about the connected users")
		return users
	}
	for data.Next() {
		user := User{}
		data.Scan(&user.Id, &user.Username)
		users = append(users, user)
	}
	return users
}
