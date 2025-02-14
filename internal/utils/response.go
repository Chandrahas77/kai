package utils

import (
	"encoding/json"
	"net/http"
)

//defines the structure for all API responses
type JSONResponse struct {
	Status  string      `json:"status"`         
	Message string      `json:"message"`        
	Data    interface{} `json:"data,omitempty"` // for optional response data
}

// this sends a structured JSON response
func RespondWithJSON(w http.ResponseWriter, statusCode int, status string, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := JSONResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// this sends an error response with a given status code
func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	RespondWithJSON(w, statusCode, "error", message, nil)
}
