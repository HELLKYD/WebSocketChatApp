package db

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id        int
	Username  string
	Password  string
	Connected bool
}

type validTypes interface {
	~string | ~int
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
