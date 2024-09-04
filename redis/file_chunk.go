package redis

import (
	"fmt"
	"strconv"

	"github.com/bladewaltz9/file-store-server/models"
)

// StoreFileInfo: store the file info in the redis
func StoreFileChunkInfo(fileInfo *models.FileChunkInfo) error {
	key := fmt.Sprintf("file_info:%s", fileInfo.FileID)

	// check if the file info exists
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to check if the file info exists: %v", err)
	}
	// if the file info exists, return directly
	if exists == 1 {
		return nil
	}

	// if the file info does not exist, store it
	_, err = rdb.HMSet(ctx, key, map[string]interface{}{
		"file_id":      fileInfo.FileID,
		"file_name":    fileInfo.FileName,
		"total_chunks": fileInfo.TotalChunks,
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to store the file info: %v", err)
	}

	return nil
}

// GetFileInfo: get the file info from the redis
func GetFileChunkInfo(fileID string) (*models.FileChunkInfo, error) {
	key := fmt.Sprintf("file_info:%s", fileID)

	// get the file info
	fileInfo, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get the file info: %v", err)
	}

	// check if the file info exists
	if len(fileInfo) == 0 {
		return nil, fmt.Errorf("file info does not exist")
	}

	totalChunks, err := strconv.Atoi(fileInfo["total_chunks"])
	if err != nil {
		return nil, fmt.Errorf("failed to convert total_chunks to int: %v", err)
	}

	// convert the file info to the structure
	return &models.FileChunkInfo{
		FileID:      fileInfo["file_id"],
		FileName:    fileInfo["file_name"],
		TotalChunks: totalChunks,
	}, nil
}

// StoreChunkStatus: store the chunk status in the redis
func StoreChunkStatus(fileID string, chunkIndex int) error {
	key := fmt.Sprintf("file_chunks:%s", fileID)

	// add the chunk index to the set
	if err := rdb.SAdd(ctx, key, chunkIndex).Err(); err != nil {
		return fmt.Errorf("failed to store the chunk status: %v", err)
	}

	return nil
}

// GetChunkStatus: get the chunk status from the redis
func GetChunkStatus(fileID string, chunkIndex int) (bool, error) {
	key := fmt.Sprintf("file_chunks:%s", fileID)

	// check if the chunk index exists
	exists, err := rdb.SIsMember(ctx, key, chunkIndex).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if the chunk index exists: %v", err)
	}

	return exists, nil
}
