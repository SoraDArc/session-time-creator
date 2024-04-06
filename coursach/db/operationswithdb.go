package db

import (
	"database/sql"
	"fmt"
)

var DbConnect *sql.DB

// Connect to DB.
func Connect(user string, password string, dbname string) {
	connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	var err error
	DbConnect, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
}

// Close connection to DB.
func Close() {
	DbConnect.Close()
}
