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
	UserID   int        `json:"user_id"`
	Username string     `json:"username"`
	Files    []FileInfo `json:"files"`
}
