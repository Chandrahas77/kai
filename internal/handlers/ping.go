package handlers

import (
	"kai-sec/internal/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func setPingRoutes(router *httprouter.Router) {
	router.GET("/ping", Ping)
}

func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	utils.RespondWithJSON(w, http.StatusOK, "success", "pong", nil)
}