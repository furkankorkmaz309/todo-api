package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
		Message: msg,
	})
}

func respondError(w http.ResponseWriter, logger *log.Logger, status int, clientMsg string, err error) {
	if logger != nil && err != nil {
		logger.Printf("[ERROR %d] %s: %v", status, clientMsg, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   clientMsg,
	})
}

func respondSuccess(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Message: message,
	})
}
