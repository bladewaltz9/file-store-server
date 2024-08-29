package db

import (
	"database/sql"
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

// UserFileExists: check if the file exists in the tbl_user_file
func UserFileExists(userID int, fileID int) (bool, error) {
	query := "SELECT id FROM tbl_user_file WHERE user_id = ? AND file_id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(userID, fileID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// SaveUserFile: save the user file relationship to the database
func SaveUserFile(userID int, fileID int, fileName string) error {
	// Begin the transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin the transaction: %v", err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Save the user file relationship
	queryInsert := "INSERT INTO tbl_user_file (user_id, file_id, file_name) VALUES (?, ?, ?)"
	stmtInsert, err := tx.Prepare(queryInsert)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmtInsert.Close()

	if _, err := stmtInsert.Exec(userID, fileID, fileName); err != nil {
		return fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	// Update the reference count
	queryUpdate := "UPDATE tbl_file SET reference_count = reference_count + 1 WHERE id = ?"
	stmtUpdate, err := tx.Prepare(queryUpdate)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmtUpdate.Close()

	if _, err := stmtUpdate.Exec(fileID); err != nil {
		return fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %v", err.Error())
	}

	return nil
}

// DeleteUserFile: delete the user file relationship from the database and delete the file if the reference count is 0
func DeleteUserFile(userID int, fileID int) (bool, string, error) {
	// Begin the transaction
	tx, err := db.Begin()
	if err != nil {
		return false, "", fmt.Errorf("failed to begin the transaction: %v", err.Error())
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete the user file relationship
	queryDelete := "DELETE FROM tbl_user_file WHERE user_id = ? AND file_id = ?"
	stmtDelete, err := tx.Prepare(queryDelete)
	if err != nil {
		return false, "", fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmtDelete.Close()

	if _, err := stmtDelete.Exec(userID, fileID); err != nil {
		return false, "", fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	// Update the reference count
	queryUpdate := "UPDATE tbl_file SET reference_count = GREATEST(reference_count - 1, 0) WHERE id = ?"
	stmtUpdate, err := tx.Prepare(queryUpdate)
	if err != nil {
		return false, "", fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmtUpdate.Close()

	if _, err := stmtUpdate.Exec(fileID); err != nil {
		return false, "", fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	// Get the reference count and file path
	querySelect := "SELECT reference_count, file_path FROM tbl_file WHERE id = ?"
	stmtSelect, err := tx.Prepare(querySelect)
	if err != nil {
		return false, "", fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmtSelect.Close()

	var referenceCount int
	var filePath string
	if err := stmtSelect.QueryRow(fileID).Scan(&referenceCount, &filePath); err != nil {
		return false, "", fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	// Delete the file if the reference count is 0
	if referenceCount == 0 {
		queryDelete := "DELETE FROM tbl_file WHERE id = ?"
		if _, err := tx.Exec(queryDelete, fileID); err != nil {
			return false, "", fmt.Errorf("failed to execute the query: %v", err.Error())
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return false, "", fmt.Errorf("failed to commit the transaction: %v", err.Error())
	}

	return referenceCount == 0, filePath, nil
}
