package handler

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/bladewaltz9/file-store-server/db"
	"github.com/bladewaltz9/file-store-server/models"
	"github.com/dgrijalva/jwt-go"
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

		// check if the user exists
		if _, err := db.GetUserInfoByUsername(username); err == nil {
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}

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

		http.Redirect(w, r, "/user/login", http.StatusFound)
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
		userInfo, err := db.GetUserInfoByUsername(username)
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

		// generate jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  userInfo.UserID,
			"username": username,
			"exp":      time.Now().Add(time.Hour * 1).Unix(), // set expiration time
		})

		tokenStr, err := token.SignedString([]byte("file-store-server"))
		if err != nil {
			log.Printf("failed to generate token: %v", err.Error())
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		// set the token in the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenStr,
			HttpOnly: true, // prevent the client from accessing the cookie
			Path:     "/",
		})

		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}
}

// DashboardHandler: handles the dashboard request
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// get the claims from the context
	claims := r.Context().Value(models.ContextKey("claims")).(jwt.MapClaims)
	user_id := int(claims["user_id"].(float64))
	username := claims["username"].(string)

	// get the user files from the database
	userFiles, err := db.GetUserFiles(user_id)
	if err != nil {
		log.Printf("failed to get user files: %v", err.Error())
		http.Error(w, "failed to get user files", http.StatusInternalServerError)
		return
	}
	data := models.DashboardData{
		UserID:   user_id,
		Username: username,
		Files:    userFiles,
	}

	tmp, err := template.ParseFiles("static/view/dashboard.html")
	if err != nil {
		log.Printf("failed to parse the template: %v", err.Error())
		http.Error(w, "failed to parse the template", http.StatusInternalServerError)
		return
	}

	if err = tmp.Execute(w, data); err != nil {
		log.Printf("failed to execute the template: %v", err.Error())
		http.Error(w, "failed to execute the template", http.StatusInternalServerError)
		return
	}

}
