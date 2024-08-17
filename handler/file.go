package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

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
			log.Printf("failed to get data from form: %v", err.Error())
			http.Error(w, "failed to get data from form", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		fileMetas := meta.FileMeta{
			FileName: header.Filename,
			FilePath: "/tmp/" + header.Filename,
		}

		// save the file to the local disk
		newFile, err := os.Create(fileMetas.FilePath)
		if err != nil {
			log.Printf("failed to create file: %v", err.Error())
			http.Error(w, "failed to create file", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		fileMetas.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			log.Printf("failed to save file: %v", err.Error())
			http.Error(w, "failed to save file", http.StatusInternalServerError)
			return
		}

		// immediately return the response to the client
		fmt.Fprintf(w, "upload file successfully")

		// calculate the hash of the file in the background
		go func() {
			fileMetas.FileHash, err = utils.CalculateSHA256(fileMetas.FilePath)
			if err != nil {
				log.Printf("failed to calculate hash: %v", err.Error())
				return
			}

			// save the file metadata to the database
			if err := db.SaveFileMeta(fileMetas.FileHash, fileMetas.FileName, fileMetas.FileSize, fileMetas.FilePath); err != nil {
				log.Printf("failed to save file metadata: %v", err.Error())
			}
		}()
	} else {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}
}

// FileQueryHandler: handles the query request
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	fileHash := r.FormValue("file_hash")
	if fileHash == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	fileMeta, err := db.GetFileMeta(fileHash)
	if err != nil {
		log.Printf("failed to get file metadata: %v", err.Error())
		http.Error(w, "failed to get file metadata", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fileMeta); err != nil {
		log.Printf("failed to encode the file metadata: %v", err.Error())
		http.Error(w, "failed to encode the file metadata", http.StatusInternalServerError)
	}
}

// FileDownloadHandler: handles the download request
func FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	// get the file hash from the request query
	fileHash := r.URL.Query().Get("file_hash")
	if fileHash == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// get the file metadata
	fileMeta, err := db.GetFileMeta(fileHash)
	if err != nil {
		log.Printf("failed to get file metadata: %v", err.Error())
		http.Error(w, "failed to get file metadata", http.StatusInternalServerError)
		return
	}

	// open the file
	file, err := os.Open(fileMeta.FilePath)
	if err != nil {
		log.Printf("failed to open file: %v", err.Error())
		http.Error(w, "failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// set the response header
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileMeta.FileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileMeta.FileSize))

	// send the file content to the client
	http.ServeContent(w, r, fileMeta.FileName, fileMeta.UpdateAt, file)
}

// FileUpdateHandler: handles the update request
func FileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	// get the file hash from the request path
	fileHash := strings.TrimPrefix(r.URL.Path, "/file/update/")
	if fileHash == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// decode the request body
	var updateReq meta.UpdateFileMetaRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		log.Printf("failed to decode the request: %v", err.Error())
		http.Error(w, "failed to decode the request", http.StatusBadRequest)
		return
	}

	// update the file metadata
	if err := db.UpdateFileMeta(fileHash, updateReq); err != nil {
		log.Printf("failed to update file metadata: %v", err.Error())
		http.Error(w, "failed to update file metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "update file metadata successfully")
}

// FileDeleteHandler: handles the delete request
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	// get the file hash from the request path
	fileHash := r.URL.Query().Get("file_hash")
	if fileHash == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// delete file from the local disk
	fileMeta, err := db.GetFileMeta(fileHash)
	if err != nil {
		log.Printf("failed to get file metadata: %v", err.Error())
		http.Error(w, "failed to get file metadata", http.StatusInternalServerError)
		return
	}

	if err := os.Remove(fileMeta.FilePath); err != nil {
		log.Printf("failed to delete file: %v", err.Error())
		http.Error(w, "failed to delete file", http.StatusInternalServerError)
		return
	}

	// delete file metadata from the database
	if err := db.DeleteFileMeta(fileHash); err != nil {
		log.Printf("failed to delete file metadata: %v", err.Error())
		http.Error(w, "failed to delete file metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "delete file successfully")
}
