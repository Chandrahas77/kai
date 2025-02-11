package handlers

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func setPingRoutes(router *httprouter.Router) {
	router.GET("/ping", Ping)
}

func Ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "pong")
}
