package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/bladewaltz9/file-store-server/db"
	"github.com/bladewaltz9/file-store-server/models"
	"github.com/bladewaltz9/file-store-server/utils"
	"github.com/google/uuid"
)

// FileChunkedUploadHandler: handles the chunked upload request
func FileChunkedUploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, "error", "invalid method")
		return
	}

	// parse the form data
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

	fileIDStr := r.FormValue("file_id")
	fileName := r.FormValue("file_name")
	chunkIndex, err := strconv.Atoi(r.FormValue("chunk_index"))
	if err != nil {
		log.Printf("failed to convert chunk_index to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	totalChunks, err := strconv.Atoi(r.FormValue("total_chunks"))
	if err != nil {
		log.Printf("failed to convert total_chunks to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}

	// get the file chunk
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("failed to get data from form: %v", err.Error())
		http.Error(w, "failed to get data from form", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// create the file directory
	if err := os.MkdirAll(filepath.Join(config.FileChunkDir, fileIDStr), os.ModePerm); err != nil {
		log.Printf("failed to create file directory: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to create file directory")
		return
	}

	// save the file chunk to the local disk
	chunkPath := filepath.Join(config.FileChunkDir, fileIDStr, fmt.Sprintf("chunk-%d", chunkIndex))
	newFile, err := os.Create(chunkPath)
	if err != nil {
		log.Printf("failed to create file: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to create file")
		return
	}
	defer newFile.Close()

	if _, err = io.Copy(newFile, file); err != nil {
		log.Printf("failed to save file: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to save file")
		return
	}

	// Record the received chunk
	val, _ := models.ChunkStatusMap.LoadOrStore(fileIDStr, &models.FileChunkInfo{
		FileID:         fileIDStr,
		FileName:       fileName,
		TotalChunks:    totalChunks,
		ReceivedChunks: make(map[int]bool),
	})
	chunkInfo := val.(*models.FileChunkInfo)
	chunkInfo.ReceivedChunks[chunkIndex] = true

	// response to the client
	utils.WriteJSONResponse(w, http.StatusOK, "success", "chunk uploaded successfully")
}

// FileChunksMergeHandler: handles the merge request
func FileChunksMergeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, "error", "invalid method")
		return
	}

	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Printf("failed to convert user_id to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	fileIDStr := r.FormValue("file_id")

	// get the file chunk info
	val, ok := models.ChunkStatusMap.Load(fileIDStr)
	if !ok {
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "file not found")
		return
	}
	chunkInfo := val.(*models.FileChunkInfo)

	// check if all chunks are received
	if len(chunkInfo.ReceivedChunks) != chunkInfo.TotalChunks {
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "not all chunks are received")
		return
	}

	// merge the file chunks
	chunkDir := filepath.Join(config.FileChunkDir, chunkInfo.FileID)
	fileMetas := &models.FileMeta{
		FileName: chunkInfo.FileName,
		FilePath: config.FileStoreDir + uuid.New().String() + "_" + chunkInfo.FileName,
		FileSize: 0,
	}

	// create the file directory
	if err := os.MkdirAll(filepath.Dir(fileMetas.FilePath), os.ModePerm); err != nil {
		log.Printf("failed to create file directory: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to create file directory")
		return
	}

	// create the new file
	newFile, err := os.Create(fileMetas.FilePath)
	if err != nil {
		log.Printf("failed to create file: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to create file")
		return
	}
	defer newFile.Close()

	// merge the file chunks
	for i := 0; i < chunkInfo.TotalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk-%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			log.Printf("failed to open chunk file: %v", err.Error())
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to open chunk file")
			return
		}
		defer chunkFile.Close()

		size, err := io.Copy(newFile, chunkFile)
		if err != nil {
			log.Printf("failed to merge chunk file: %v", err.Error())
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to merge chunk file")
			return
		}
		fileMetas.FileSize += size
	}

	// delete the file chunks
	if err := os.RemoveAll(chunkDir); err != nil {
		log.Printf("failed to delete chunk directory: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to delete chunk directory")
		return
	}

	// calculate the hash of the file
	fileMetas.FileHash, err = utils.CalculateSHA256(newFile)
	if err != nil {
		log.Printf("failed to calculate hash: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to calculate hash")
		return
	}

	// save the file metadata to the database
	exist, fileID, err := db.FileExists(fileMetas.FileHash)
	if err != nil {
		log.Printf("failed to check if the file exists: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to check if the file exists")
		return
	}
	if !exist { // If the file does not exist, save the file metadata to the database
		fileID, err = db.SaveFileMeta(fileMetas.FileHash, fileMetas.FileName, fileMetas.FileSize, fileMetas.FilePath)
		if err != nil {
			log.Printf("failed to save file metadata: %v", err.Error())
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to save file metadata")
			return
		}
	} else { // If the file exists, delete the file
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
		if err := db.SaveUserFile(userID, fileID, fileMetas.FileName); err != nil {
			log.Printf("failed to save user file: %v", err.Error())
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to save user file")
			return
		}
	} else {
		utils.WriteJSONResponse(w, http.StatusOK, "success", "file already exists")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, "success", "file uploaded successfully")
}
