package utils

import (
	"encoding/json"
	"net/http"
)

// JSONResponse represents a standard JSON response structure
type JSONResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// writeJSONResponse writes a JSON response to the client
func WriteJSONResponse(w http.ResponseWriter, statusCode int, status, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := JSONResponse{
		Status:  status,
		Message: message,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"status": "error", "message": "failed to encode JSON response"}`, http.StatusInternalServerError)
	}
}
