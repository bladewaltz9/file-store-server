package meta

import "time"

// FileMeta: file metadata structure
type FileMeta struct {
	FileHash   string    `json:"file_hash"`
	FileName   string    `json:"file_name"`
	FileSize   int64     `json:"file_size"`
	FilePath   string    `json:"file_path"`
	UploadTime time.Time `json:"upload_time"`
	Status     string    `json:"status"`
}

// UpdateFileMetaReq: update file metadata request structure
type UpdateFileMetaReq struct {
	FileName string `json:"file_name"`
	Status   string `json:"status"`
}

// var fileMetas map[string]FileMeta
// var mu sync.Mutex

// func init() {
// 	fileMetas = make(map[string]FileMeta)
// }

// // AddFileMeta: add the file metadata
// func AddFileMeta(fileMeta FileMeta) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	fileMetas[fileMeta.FileHash] = fileMeta
// }

// // GetFileMeta: get the file metadata by file hash
// func GetFileMeta(fileHash string) (FileMeta, error) {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if meta, ok := fileMetas[fileHash]; ok {
// 		return meta, nil
// 	}
// 	return FileMeta{}, errors.New("file not found")
// }

// // UpladFileMeta: update the file metadata
// func UpdateFileMeta(fileHash string, fileMeta FileMeta) error {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if _, ok := fileMetas[fileHash]; ok {
// 		fileMetas[fileHash] = fileMeta
// 		return nil
// 	}
// 	return errors.New("file not found")
// }

// // DeleteFileMeta: delete the file metadata
// func DeleteFileMeta(fileHash string) error {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	if _, ok := fileMetas[fileHash]; ok {
// 		delete(fileMetas, fileHash)
// 		return nil
// 	}
// 	return errors.New("file not found")
// }

// // ListAllFileMetas: list all the file metadata
// func ListAllFileMetas() []FileMeta {
// 	mu.Lock()
// 	defer mu.Unlock()
// 	metas := make([]FileMeta, 0, len(fileMetas))
// 	for _, meta := range fileMetas {
// 		metas = append(metas, meta)
// 	}
// 	return metas
// }
