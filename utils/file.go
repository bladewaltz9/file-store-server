package utils

import (
	"fmt"

	"github.com/bladewaltz9/file-store-server/db"
	"github.com/bladewaltz9/file-store-server/models"
)

// SaveUserFileDB saves the file metadata to the database
func SaveUserFileDB(fileMetas *models.FileMeta, userID int) error {
	// save the file metadata to the database
	fileID, err := db.SaveFileMeta(fileMetas.FileHash, fileMetas.FileName, fileMetas.FileSize, fileMetas.FilePath)
	if err != nil {
		return fmt.Errorf("failed to save file metadata: %v", err.Error())
	}

	// save the relationship between the user and the file to the database
	if err := db.SaveUserFile(userID, fileID, fileMetas.FileName); err != nil {
		return fmt.Errorf("failed to save user file: %v", err.Error())
	}

	return nil
}
