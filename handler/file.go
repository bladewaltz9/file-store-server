package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bladewaltz9/file-store-server/db"
	"github.com/bladewaltz9/file-store-server/meta"
	"github.com/bladewaltz9/file-store-server/utils"
)

// FileUploadHandler: handles the upload request
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// return the upload page
		http.ServeFile(w, r, "static/view/file_upload.html")
	} else if r.Method == http.MethodPost {
		// handle the upload file
		r.ParseMultipartForm(32 << 20) // limit the file size to 32MB
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Failed to get data from form: %v", err.Error())
			http.Error(w, "Failed to get data from form", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileMetas := meta.FileMeta{
			FileName:   header.Filename,
			FilePath:   "/tmp/" + header.Filename,
			UploadTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		// save the file to the local disk
		newFile, err := os.Create(fileMetas.FilePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err.Error())
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		fileMetas.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			log.Printf("Failed to save file: %v", err.Error())
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// immediately return the response to the client
		fmt.Fprintf(w, "Upload file successfully")

		// calculate the hash of the file in the background
		go func() {
			fileMetas.FileHash, err = utils.CalculateSHA256(fileMetas.FilePath)
			if err != nil {
				log.Printf("Failed to calculate hash: %v", err.Error())
				return
			}

			// save the file metadata to the database
			if err := db.SaveFileMeta(fileMetas.FileHash, fileMetas.FileName, fileMetas.FileSize, fileMetas.FilePath); err != nil {
				log.Printf("Failed to save file metadata: %v", err.Error())
			}
		}()
	}
}

// FileQueryHandler: handles the query request
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
}

// FileDownloadHandler: handles the download request
func FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
}

// FileUpdateHandler: handles the update request
func FileUpdateHandler(w http.ResponseWriter, r *http.Request) {
}

// FileDeleteHandler: handles the delete request
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
}
