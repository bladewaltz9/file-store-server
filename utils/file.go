package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/bladewaltz9/file-store-server/db"
	"github.com/bladewaltz9/file-store-server/models"
)

// SaveUserFileDB saves the file metadata to the database
func SaveUserFileDB(fileMetas *models.FileMeta, userID int) error {
	// save the file metadata to the database
	exist, fileID, err := db.FileExists(fileMetas.FileHash)
	if err != nil {
		return fmt.Errorf("failed to check if the file exists: %v", err.Error())
	}
	if !exist {
		// If the file does not exist, save the file metadata to the database
		fileID, err = db.SaveFileMeta(fileMetas.FileHash, fileMetas.FileName, fileMetas.FileSize, fileMetas.FilePath)
		if err != nil {
			return fmt.Errorf("failed to save file metadata: %v", err.Error())
		}
	} else { // If the file exists, delete the file
		go func() {
			if err := os.Remove(fileMetas.FilePath); err != nil {
				log.Printf("failed to delete file: %v", err.Error())
			}
		}()
	}

	// check if the file exists in the user file table
	exist, err = db.UserFileExists(userID, fileID)
	if err != nil {
		return fmt.Errorf("failed to check if the file exists: %v", err.Error())
	}
	if !exist {
		if err := db.SaveUserFile(userID, fileID, fileMetas.FileName); err != nil {
			return fmt.Errorf("failed to save user file: %v", err.Error())
		}
	} else {
		return fmt.Errorf("file already exists")
	}

	return nil
}
