package main

import (
	"net/http"

	"github.com/bladewaltz9/file-store-server/handler"
)

func main() {
	http.HandleFunc("/upload", handler.UploadHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
