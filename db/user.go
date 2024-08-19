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
func GetUserInfoByUsername(username string) (*models.UserInfo, error) {
	query := "SELECT id, username, password, email FROM tbl_user WHERE username = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	user := &models.UserInfo{}
	err = stmt.QueryRow(username).Scan(&user.UserID, &user.Username, &user.Password, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	return user, nil
}

// GetUserFiles: get the user files from the database
func GetUserFiles(user_id int) ([]models.FileInfo, error) {
	query := `SELECT f.id, f.file_name, f.file_size, DATE_FORMAT(uf.upload_at, '%Y-%m-%d %H:%i'), uf.status 
	FROM tbl_user_file uf
	JOIN tbl_file f ON uf.file_id = f.id 
	WHERE uf.user_id = ?;`

	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	defer rows.Close()

	var userFiles []models.FileInfo
	for rows.Next() {
		file := models.FileInfo{}
		if err := rows.Scan(&file.FileID, &file.FileName, &file.FileSize, &file.UploadTime, &file.Status); err != nil {
			return nil, fmt.Errorf("failed to scan the row: %v", err.Error())
		}
		userFiles = append(userFiles, file)
	}
	return userFiles, nil
}
