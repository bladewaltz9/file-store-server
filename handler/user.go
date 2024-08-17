package handler

import (
	"log"
	"net/http"

	"github.com/bladewaltz9/file-store-server/db"
	"golang.org/x/crypto/bcrypt"
)

// UserRegisterHandler: handles the user register request
func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "static/view/user_register.html")
	} else if r.Method == http.MethodPost {
		// parse the form
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		// encode the password
		encodedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("failed to encode the password: %v", err.Error())
			http.Error(w, "failed to encode the password", http.StatusInternalServerError)
			return
		}

		// save user to the database
		if err := db.SaveUserInfo(username, string(encodedPwd), email); err != nil {
			log.Printf("failed to save user: %v", err.Error())
			http.Error(w, "failed to save user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("register successfully"))
	} else {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}
}

// UserLoginHandler: handles the user login request
func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "static/view/user_login.html")
	} else if r.Method == http.MethodPost {
		// parse the form
		username := r.FormValue("username")
		password := r.FormValue("password")

		// get the user information from the database
		userInfo, err := db.GetUserInfo(username)
		if err != nil {
			log.Printf("failed to get user: %v", err.Error())
			http.Error(w, "failed to get user", http.StatusInternalServerError)
			return
		}

		// compare the password
		if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(password)); err != nil {
			log.Printf("invalid password: %v", err.Error())
			http.Error(w, "invalid password", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("login successfully"))
	} else {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}
}
