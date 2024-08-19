package models

// FileInfo: file information structure in dashboard
type FileInfo struct {
	FileID     int    `json:"file_id"`
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	UploadTime string `json:"upload_time"`
	Status     string `json:"status"`
}

// DashboardData: dashboard data structure
type DashboardData struct {
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Files    []FileInfo `json:"files"`
}
