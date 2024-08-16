package db

import (
	"fmt"

	"github.com/bladewaltz9/file-store-server/meta"
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
func GetFileMeta(fileHash string) (*meta.FileMeta, error) {
	query := "SELECT file_hash, file_name, file_size, file_path, create_at FROM tbl_file WHERE file_hash = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	fileMeta := &meta.FileMeta{}
	err = stmt.QueryRow(fileHash).Scan(&fileMeta.FileHash, &fileMeta.FileName, &fileMeta.FileSize, &fileMeta.FilePath, &fileMeta.UploadTime)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %v", err.Error())
	}
	return fileMeta, nil
}

// UpdateFileMeta: update the file metadata in the database
func UpdateFileMeta(fileHash string, updateReq meta.UpdateFileMetaReq) error {
	query := "UPDATE tbl_file SET file_name = ?, status = ? WHERE file_hash = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(updateReq.FileName, updateReq.Status, fileHash)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %v", err.Error())
	}

	return nil
}
