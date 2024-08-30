package main

import (
	"net/http"

	"github.com/bladewaltz9/file-store-server/config"
	"github.com/bladewaltz9/file-store-server/handler"
	"github.com/bladewaltz9/file-store-server/middleware"
)

func main() {
	// file handler
	http.HandleFunc("/file/upload", middleware.TokenAuthMiddleware(handler.FileUploadHandler))
	http.HandleFunc("/file/query", middleware.TokenAuthMiddleware(handler.FileQueryHandler))
	http.HandleFunc("/file/download/", middleware.TokenAuthMiddleware(handler.FileDownloadHandler))
	http.HandleFunc("/file/update/", middleware.TokenAuthMiddleware(handler.FileUpdateHandler))
	http.HandleFunc("/file/delete/", middleware.TokenAuthMiddleware(handler.FileDeleteHandler))

	// file chunked handler
	http.HandleFunc("/file/upload/chunk", middleware.TokenAuthMiddleware(handler.FileChunkedUploadHandler))
	http.HandleFunc("/file/merge", middleware.TokenAuthMiddleware(handler.FileChunksMergeHandler))

	// user handler
	http.HandleFunc("/user/register", handler.UserRegisterHandler)
	http.HandleFunc("/user/login", handler.UserLoginHandler)

	// dashboard handler
	http.HandleFunc("/dashboard", middleware.TokenAuthMiddleware(handler.DashboardHandler))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if middleware.IsAuthenticated(r) {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		} else {
			http.Redirect(w, r, "/user/login", http.StatusFound)
		}
	})

	// start the server
	err := http.ListenAndServeTLS(":8080", config.CertFile, config.KeyFile, nil)
	if err != nil {
		panic(err)
	}
}
