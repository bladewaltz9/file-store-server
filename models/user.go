package models

// UserInfo: user information structure
type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	// TODO: add more fields
}
