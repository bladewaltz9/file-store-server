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
	if err := r.ParseMultipartForm(config.MaxUploadSize); err != nil {
		if err == http.ErrContentLength {
			log.Printf("uploaded file is too large: %v", err)
			utils.WriteJSONResponse(w, http.StatusRequestEntityTooLarge, "error", "uploaded file is too large")
		} else {
			log.Printf("failed to parse multipart form: %v", err)
			utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "failed to parse form data")
		}
		return
	}

	// get the user_id, file_hash, and file from the form
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Printf("failed to convert user_id to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	fileHash := r.FormValue("file_hash")
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("failed to get data from form: %v", err.Error())
		http.Error(w, "failed to get data from form", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileMetas := &models.FileMeta{}
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

	// check the file hash with the hash from the client
	if fileMetas.FileHash != fileHash {
		// delete the file from the local disk
		go func() {
			if err := os.Remove(fileMetas.FilePath); err != nil {
				log.Printf("failed to delete file: %v", err.Error())
			}
		}()
		log.Printf("file hash does not match")
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "file hash does not match")
		return
	}

	// save the file metadata to the database
	if err := SaveUserFileDB(fileMetas, userID); err != nil {
		log.Printf("failed to save file metadata: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", err.Error())
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
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, "error", "invalid method")
		return
	}

	// get the user_id and file_id from the request path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	userID, err := strconv.Atoi(parts[3])
	if err != nil {
		log.Printf("failed to convert user_id to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	fileID, err := strconv.Atoi(parts[4])
	if err != nil {
		log.Printf("failed to convert file_id to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}

	// delete the file
	ok, filePath, err := db.DeleteUserFile(userID, fileID)
	if err != nil {
		log.Printf("failed to delete file: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	// delete the file from the local disk if the reference count is 0
	if ok {
		go func() {
			if err := os.Remove(filePath); err != nil {
				log.Printf("failed to delete file: %v", err.Error())
			}
		}()
	}

	utils.WriteJSONResponse(w, http.StatusOK, "success", "file deleted successfully")
}

// FileFastUploadHandler: fast upload the file if the file already exists
func FileFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, "error", "invalid method")
		return
	}

	// get the user_id, file_hash, and file_name from the form
	fileHash := r.FormValue("file_hash")
	fileName := r.FormValue("file_name")
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Printf("failed to convert user_id to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	if fileHash == "" || fileName == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}

	// check if the file exists
	exist, fileID, err := db.FileExists(fileHash)
	if err != nil {
		log.Printf("failed to check if the file exists: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to check if the file exists")
		return
	}

	// if the file does not exist, return the status of "not_exists"
	if !exist {
		utils.WriteJSONResponse(w, http.StatusOK, "not_exists", "file does not exist")
		return
	}

	// check if the file exists in the user file table
	exist, err = db.UserFileExists(userID, fileID)
	if err != nil {
		log.Printf("failed to check if the file exists: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to check if the file exists")
		return
	}

	// if the file exists in the user file table, return the status of "repeat"
	if exist {
		utils.WriteJSONResponse(w, http.StatusOK, "repeat", "file already exists")
		return
	}

	// save the file to the user file table
	if err := db.SaveUserFile(userID, fileID, fileName); err != nil {
		log.Printf("failed to save user file: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to save user file")
		return
	}

	// return the status of "success"
	utils.WriteJSONResponse(w, http.StatusOK, "success", "file fast uploaded successfully")
}
