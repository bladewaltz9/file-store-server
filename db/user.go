package db

import "fmt"

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
