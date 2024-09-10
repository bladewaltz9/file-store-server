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
	"github.com/bladewaltz9/file-store-server/models"
	"github.com/bladewaltz9/file-store-server/mq"
	"github.com/bladewaltz9/file-store-server/redis"
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
	chunkHash := r.FormValue("chunk_hash")
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

	// calculate the hash of the chunk
	chunkHashCalculated, err := utils.CalculateSHA256(newFile)
	if err != nil {
		log.Printf("failed to calculate hash: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to calculate hash")
		return
	}
	// check the chunk hash with the hash from the client
	if chunkHash != chunkHashCalculated {
		// delete the chunk from the local disk
		go func() {
			if err := os.Remove(chunkPath); err != nil {
				log.Printf("failed to delete chunk: %v", err.Error())
			}
		}()
		log.Printf("chunk hash does not match: %v", chunkHashCalculated)
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "chunk hash does not match")
		return
	}

	// store the file info
	chunkInfo := &models.FileChunkInfo{
		FileID:      fileIDStr,
		FileName:    fileName,
		TotalChunks: totalChunks,
	}
	if err := redis.StoreFileChunkInfo(chunkInfo); err != nil {
		log.Printf("failed to store file info: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to store file info")
		return
	}

	// store the chunk status
	if err := redis.StoreChunkStatus(fileIDStr, chunkIndex); err != nil {
		log.Printf("failed to store chunk status: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to store chunk status")
		return
	}

	// response to the client
	utils.WriteJSONResponse(w, http.StatusOK, "success", "chunk uploaded successfully")
}

// FileChunksMergeHandler: handles the merge request
func FileChunksMergeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, "error", "invalid method")
		return
	}

	// parse the form data
	userID, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		log.Printf("failed to convert user_id to int: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "invalid parameter")
		return
	}
	fileIDStr := r.FormValue("file_id")
	fileHash := r.FormValue("file_hash")

	// check if all chunks are received
	chunkInfo, err := redis.GetFileChunkInfo(fileIDStr)
	if err != nil || chunkInfo == nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "file not found")
		return
	}
	for i := 0; i < chunkInfo.TotalChunks; i++ {
		received, err := redis.GetChunkStatus(fileIDStr, i)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to get chunk status")
			return
		}
		if !received {
			utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "not all chunks are received, lost chunk: "+strconv.Itoa(i))
			return
		}
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
	go func() {
		if err := os.RemoveAll(chunkDir); err != nil {
			log.Printf("failed to delete chunk directory: %v", err.Error())
		}
	}()

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
		log.Printf("file hash does not match: %v", fileMetas.FileHash)
		utils.WriteJSONResponse(w, http.StatusBadRequest, "error", "file hash does not match")
		return
	}

	// save the file metadata to the database
	if err := SaveUserFileDB(fileMetas, userID); err != nil {
		log.Printf("failed to save file metadata: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	// send the message to the MQ
	rabbitMQ := mq.GetRabbitMQ()
	fileMsg := &mq.FileTransferMessage{
		FileID:    fileMetas.FileID,
		LocalFile: fileMetas.FilePath,
		ObjectKey: config.BucketDir + fileMetas.FileName,
	}
	if err := rabbitMQ.PublishMessage(fileMsg); err != nil {
		log.Printf("failed to publish message: %v", err.Error())
		utils.WriteJSONResponse(w, http.StatusInternalServerError, "error", "failed to publish message")
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, "success", "file uploaded successfully")
}
