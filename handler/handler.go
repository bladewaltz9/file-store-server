package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// UploadHandler: handles the upload request
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// return the upload page
		http.ServeFile(w, r, "static/view/index.html")
	} else if r.Method == http.MethodPost {
		// handle the upload file
		// r.ParseMultipartForm(32 << 20) // limit the file size to 32MB

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get data from form", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		location := "/tmp/" + header.Filename

		// save the file to the local disk
		newFile, err := os.Create(location)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		if _, err = io.Copy(newFile, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Successfully uploaded file")

	}
}
