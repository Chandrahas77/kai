package handlers

import (
	"encoding/json"
	"kai-sec/internal/services"
	"kai-sec/internal/utils"
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if req.Repository == "" || len(req.Files) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Repo and files are required")
		return
	}

	// Process Scan
	err = services.ProcessScan(req.Repository, req.Files)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process scan")
		return
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "Scan completed successfully"})
		w.WriteHeader(http.StatusOK)
	}

}
