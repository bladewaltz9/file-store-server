package db

import "fmt"

// SaveFileMeta: save the file metadata to the database
func SaveFileMeta(fileHash string, fileName string, fileSize int64, filePath string) error {
	query := "INSERT INTO tbl_file (file_hash, file_name, file_size, file_path) VALUES (?, ?, ?, ?)"

	stmt, err := db.Prepare(query)
	if err != nil {
		return fmt.Errorf("Failed to prepare the query: %v", err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileHash, fileName, fileSize, filePath)
	if err != nil {
		return fmt.Errorf("Failed to execute the query: %v", err.Error())
	}
	return nil
}
