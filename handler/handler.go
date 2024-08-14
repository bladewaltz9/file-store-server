package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bladewaltz9/file-store-server/meta"
	"github.com/bladewaltz9/file-store-server/utils"
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

		fileMetas := meta.FileMeta{
			FileName:   header.Filename,
			FilePath:   "/tmp/" + header.Filename,
			UploadTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		// save the file to the local disk
		newFile, err := os.Create(fileMetas.FilePath)
		if err != nil {
			http.Error(w, "Failed to create file", http.StatusInternalServerError)
			return
		}
		defer newFile.Close()

		if _, err = io.Copy(newFile, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		// calculate the hash of the file
		// TODO: 计算大文件的 hash 会比较久，可以考虑异步计算
		fileMetas.FileHash, err = utils.CalculateSHA256(newFile)
		if err != nil {
			http.Error(w, "Failed to calculate hash", http.StatusInternalServerError)
			return
		}

		// save the file metadata
		meta.AddFileMeta(fileMetas)

		fmt.Println(fileMetas)

		fmt.Fprintln(w, "Successfully uploaded file")

	}
}
