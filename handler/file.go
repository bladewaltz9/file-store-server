package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/bladewaltz9/file-store-server/db"
	"github.com/bladewaltz9/file-store-server/models"
	"github.com/bladewaltz9/file-store-server/utils"
	"github.com/google/uuid"
)

// FileUploadHandler: handles the upload request
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, "error", "invalid method")
		return
	}

	// handle the upload file
	r.ParseMultipartForm(config.MaxUploadSize) // limit the file size

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Printf("failed to convert user_id to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("failed to get data from form: %v", err.Error())
		http.Error(w, "failed to get data from form", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileMetas := models.FileMeta{}
	fileMetas.FileName = header.Filename
	fileMetas.FilePath = config.FileStoreDir + uuid.New().String() + "_" + header.Filename

	// create the file directory
	fileDir := filepath.Dir(fileMetas.FilePath)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		log.Printf("failed to create file directory: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to create file directory")
		return
	}

	// save the file to the local disk
	newFile, err := os.Create(fileMetas.FilePath)
	if err != nil {
		log.Printf("failed to create file: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to create file")
		return
	}
	defer newFile.Close()

	fileMetas.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		log.Printf("failed to save file: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to save file")
		return
	}

	// calculate the hash of the file
	fileMetas.FileHash, err = utils.CalculateSHA256(newFile)
	if err != nil {
		log.Printf("failed to calculate hash: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to calculate hash")
		return
	}

	// check if the file exists in the file table
	exist, fileID, err := db.FileExists(fileMetas.FileHash)
	if err != nil {
		log.Printf("failed to check if the file exists: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to check if the file exists")
		return
	}
	if !exist {
		// save the file metadata to the database
		fileID, err = db.SaveFileMeta(fileMetas.FileHash, fileMetas.FileName, fileMetas.FileSize, fileMetas.FilePath)
		if err != nil {
			log.Printf("failed to save file metadata: %v", err.Error())
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to save file metadata")
			return
		}
	} else {
		if err := os.Remove(fileMetas.FilePath); err != nil {
			log.Printf("failed to delete file: %v", err.Error())
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to delete file")
			return
		}
	}

	// check if the file exists in the user file table
	exist, err = db.UserFileExists(userID, fileID)
	if err != nil {
		log.Printf("failed to check if the file exists: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to check if the file exists")
		return
	}
	if !exist {
		// save the user file relationship to the database
		if err := db.SaveUserFile(userID, fileID, fileMetas.FileName); err != nil {
			log.Printf("failed to save user file: %v", err.Error())
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to save user file")
			return
		}
	} else {
		utils.WriteJSONResponse(w, http.StatusOK, "error", "file already exists")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, "success", "file uploaded successfully")
}

// FileQueryHandler: handles the query request
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}

	fileIDStr := r.FormValue("file_id")
	if fileIDStr == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// convert the file_id to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		log.Printf("failed to convert file_id to int: %v", err.Error())
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	fileMeta, err := db.GetFileMeta(fileID)
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

	fileIDStr := strings.TrimPrefix(r.URL.Path, "/file/download/")
	if fileIDStr == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// convert the file_id to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		log.Printf("failed to convert file_id to int: %v", err.Error())
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// get the file metadata
	fileMeta, err := db.GetFileMeta(fileID)
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
	fileIDStr := strings.TrimPrefix(r.URL.Path, "/file/update/")
	if fileIDStr == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// convert the file_id to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		log.Printf("failed to convert file_id to int: %v", err.Error())
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// decode the request body
	var updateReq models.UpdateFileMetaRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		log.Printf("failed to decode the request: %v", err.Error())
		http.Error(w, "failed to decode the request", http.StatusBadRequest)
		return
	}

	// update the file metadata
	if err := db.UpdateFileMeta(fileID, updateReq); err != nil {
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
	fileIDStr := r.URL.Query().Get("file_id")
	if fileIDStr == "" {
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// convert the file_id to int
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil {
		log.Printf("failed to convert file_id to int: %v", err.Error())
		http.Error(w, "invalid parameter", http.StatusBadRequest)
		return
	}

	// delete file from the local disk
	fileMeta, err := db.GetFileMeta(fileID)
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
	if err := db.DeleteFileMeta(fileID); err != nil {
		log.Printf("failed to delete file metadata: %v", err.Error())
		http.Error(w, "failed to delete file metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "delete file successfully")
}
