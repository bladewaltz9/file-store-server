package db

import (
	"fmt"

	"github.com/bladewaltz9/file-store-server/models"
)

// SaveFileMeta: save the file metadata to the database
func SaveFileMeta(fileHash string, fileName string, fileSize int64, filePath string) error {
	query := "INSERT INTO tbl_file (file_hash, file_name, file_size, file_path) VALUES (?, ?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileHash, fileName, fileSize, filePath)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	return nil
}

// GetFileMeta: get the file metadata from the database
func GetFileMeta(fileID string) (*models.FileMeta, error) {
	query := "SELECT file_hash, file_name, file_size, file_path, create_at, update_at, status FROM tbl_file WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	fileMeta := &models.FileMeta{}
	err = stmt.QueryRow(fileID).Scan(&fileMeta.FileHash, &fileMeta.FileName, &fileMeta.FileSize, &fileMeta.FilePath, &fileMeta.CreateAt, &fileMeta.UpdateAt, &fileMeta.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	return fileMeta, nil
}

// UpdateFileMeta: update the file metadata in the database
func UpdateFileMeta(fileID string, updateReq models.UpdateFileMetaRequest) error {
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

// DeleteFileMeta: delete the file metadata from the database
func DeleteFileMeta(fileID string) error {
	query := "DELETE FROM tbl_file WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileID)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	return nil
}
