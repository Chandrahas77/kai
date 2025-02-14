package main

import (
	"kai-sec/internal/handlers"
	"kai-sec/internal/logger"
	"kai-sec/internal/storage"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	logger.InitLogger()
	l := logger.Log
	db, err := storage.ConnectDB()
	if err != nil {
		l.Fatal("Database connection failed:", zap.Error(err))
	}
	defer db.Close()

	l.Info("Service started successfully!")

	router := handlers.NewRouter()

	// Start the server
	//TODO make it env variable
	port := "8080"
	l.Info("Server running on port", zap.String("port", port))
	http.ListenAndServe(":"+port, router)
}
