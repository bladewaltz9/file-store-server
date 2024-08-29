package db

import (
	"database/sql"
	"fmt"

	"github.com/bladewaltz9/file-store-server/models"
)

// SaveFileMeta: save the file metadata to the database
func SaveFileMeta(fileHash string, fileName string, fileSize int64, filePath string) (int, error) {
	query := "INSERT INTO tbl_file (file_hash, file_name, file_size, file_path) VALUES (?, ?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	result, err := stmt.Exec(fileHash, fileName, fileSize, filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	fileID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get the last insert id: %v", err.Error())
	}

	return int(fileID), nil
}

// GetFileMeta: get the file metadata from the database
func GetFileMeta(fileID int) (*models.FileMeta, error) {
	query := "SELECT file_hash, file_name, file_size, file_path, create_at, update_at, status FROM tbl_file WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	fileMeta := &models.FileMeta{
		FileID: fileID,
	}
	err = stmt.QueryRow(fileID).Scan(&fileMeta.FileHash, &fileMeta.FileName, &fileMeta.FileSize, &fileMeta.FilePath, &fileMeta.CreateAt, &fileMeta.UpdateAt, &fileMeta.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	return fileMeta, nil
}

// UpdateFileMeta: update the file metadata in the database
func UpdateFileMeta(fileID int, updateReq models.UpdateFileMetaRequest) error {
	query := "UPDATE tbl_file SET file_name = ?, status = ? WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(updateReq.FileName, updateReq.Status, fileID)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	return nil
}

// // DeleteFileMeta: delete the file metadata from the database
// func DeleteFileMeta(fileID int) error {
// 	query := "DELETE FROM tbl_file WHERE id = ?"

// 	stmt, err := db.Prepare(query)
// 	if err != nil {
// 		return fmt.Errorf("failed to prepare the query: %v", err.Error())
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(fileID)
// 	if err != nil {
// 		return fmt.Errorf("failed to execute the query: %v", err.Error())
// 	}

// 	return nil
// }

// FileExists: check if the file exists in the tbl_file
func FileExists(fileHash string) (bool, int, error) {
	query := "SELECT id FROM tbl_file WHERE file_hash = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return false, 0, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	var fileID int
	err = stmt.QueryRow(fileHash).Scan(&fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, err
	}

	return true, fileID, nil
}
