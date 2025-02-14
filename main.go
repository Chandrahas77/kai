package main

import (
	"kai-sec/internal/config"
	"kai-sec/internal/handlers"
	"kai-sec/internal/logger"
	"kai-sec/internal/storage"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	// Load configuration
	logger.InitLogger()
	l := logger.GetLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal("Database connection failed:", zap.Error(err))
	}

	// Initialize database connection
	db, err := storage.InitDB(cfg)
	if err != nil {
		l.Fatal("Could not connect to database: %v", zap.Error(err))
	}
	defer db.Close()

	router := handlers.NewRouter()

	// Start the server
	l.Info("Server running on port", zap.String("port", cfg.ServerPort))
	http.ListenAndServe(":"+cfg.ServerPort, router)
}
