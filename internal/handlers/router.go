package handlers

import (
	"github.com/julienschmidt/httprouter"
)

// Initialize and return the router
func NewRouter() *httprouter.Router {
	router := httprouter.New()
	setPingRoutes(router)
	return router
}
