package main

import (
	"net/http"

	"github.com/bladewaltz9/file-store-server/handler"
	"github.com/bladewaltz9/file-store-server/middleware"
)

func main() {
	// file handler
	http.HandleFunc("/file/upload", middleware.TokenAuthMiddleware(handler.FileUploadHandler))
	http.HandleFunc("/file/query", middleware.TokenAuthMiddleware(handler.FileQueryHandler))
	http.HandleFunc("/file/download/", middleware.TokenAuthMiddleware(handler.FileDownloadHandler))
	http.HandleFunc("/file/update/", middleware.TokenAuthMiddleware(handler.FileUpdateHandler))
	http.HandleFunc("/file/delete", middleware.TokenAuthMiddleware(handler.FileDeleteHandler))

	// user handler
	http.HandleFunc("/user/register", handler.UserRegisterHandler)
	http.HandleFunc("/user/login", handler.UserLoginHandler)

	http.HandleFunc("/dashboard", middleware.TokenAuthMiddleware(handler.DashboardHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
