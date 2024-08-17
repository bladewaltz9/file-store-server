package db

import (
	"fmt"

	"github.com/bladewaltz9/file-store-server/models"
)

// SaveUserInfo: save the user information to the database
func SaveUserInfo(username string, password string, email string) error {
	query := "INSERT INTO tbl_user (username, password, email) VALUES (?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, password, email)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	return nil
}

// GetUserInfo: get the user information from the database
func GetUserInfo(username string) (*models.UserInfo, error) {
	query := "SELECT username, password, email FROM tbl_user WHERE username = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	user := &models.UserInfo{}
	err = stmt.QueryRow(username).Scan(&user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	return user, nil
}
