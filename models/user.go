package models

// UserInfo: user information structure
type UserInfo struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	// TODO: add more fields
}

type ContextKey string // ContextKey: context key type
