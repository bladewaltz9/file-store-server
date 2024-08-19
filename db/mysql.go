package db

import (
	"database/sql"
	"fmt"

	"github.com/bladewaltz9/file-store-server/config"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// init: initialize the mysql connection
func init() {
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to the mysql: %v", err.Error()))
	}

	db.SetMaxOpenConns(config.DBMaxConn)

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping the mysql: %v", err.Error()))
	}
}
