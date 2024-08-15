package meta

import "time"

// FileMeta: file metadata structure
type FileMeta struct {
	FileHash   string
	FileName   string
	FileSize   int64
	FilePath   string
	UploadTime time.Time
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
