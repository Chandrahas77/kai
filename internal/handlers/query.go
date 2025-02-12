package handlers

import (
	"encoding/json"
	"kai-sec/internal/dtos"
	"kai-sec/internal/services"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func setQueryRoutes(router *httprouter.Router) {
	router.POST("/query", QueryHandler)
}

func QueryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req dtos.FliterRequest
	// Decode request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Get vulnerabilities by severity
	results, err := services.GetVulnerabilities(req)
	if err != nil {
		http.Error(w, "Failed to fetch vulnerabilities", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
