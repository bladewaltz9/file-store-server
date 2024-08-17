package models

import "time"

// FileMeta: file metadata structure
type FileMeta struct {
	FileHash string    `json:"file_hash"`
	FileName string    `json:"file_name"`
	FileSize int64     `json:"file_size"`
	FilePath string    `json:"file_path"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
	Status   string    `json:"status"`
}

// UpdateFileMetaRequest: update file metadata request structure
type UpdateFileMetaRequest struct {
	FileName string `json:"file_name"`
	Status   string `json:"status"`
}
