package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// init: initialize the mysql connection
func init() {
	var err error
	db, err = sql.Open("mysql", "root:Lollzp1999!@tcp(127.0.0.1:3306)/file_server?charset=utf8")
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the mysql: %v", err.Error()))
	}

	db.SetMaxOpenConns(1000)

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping the mysql: %v", err.Error()))
	}
}
