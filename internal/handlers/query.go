package handlers

import (
	"encoding/json"
	"kai-sec/internal/dtos"
	"kai-sec/internal/services"
	"kai-sec/internal/utils"
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
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Get vulnerabilities by severity
	results, err := services.GetVulnerabilities(req)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch vulnerabilities")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.RespondWithJSON(w, http.StatusOK, "success", "Vulnerabilities fetched successfully", results)
}
