package handlers

import (
	"encoding/json"
	"kai-sec/internal/services"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ScanRequest struct {
	Repository string   `json:"repo"`
	Files      []string `json:"files"`
}

func setScanRoutes(router *httprouter.Router) {
	router.POST("/scan", ScanHandler)
}

func ScanHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req ScanRequest

	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Repository == "" || len(req.Files) == 0 {
		http.Error(w, "Repo and files are required", http.StatusBadRequest)
		return
	}

	// Process Scan
	err = services.ProcessScan(req.Repository, req.Files)
	if err != nil {
		http.Error(w, "Failed to process scan", http.StatusInternalServerError)
		return
	}

	// Respond
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Scan completed successfully"})
}
