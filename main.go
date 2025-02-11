package main

import (
	"fmt"
	"kai-sec/internal/handlers"
	"kai-sec/internal/storage"
	"log"
	"net/http"
)

func main() {
	db, err := storage.ConnectDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	log.Println("Service started successfully!")

	router := handlers.NewRouter()

	// Start the server
	//TODO make it env variable
	port := "8080"
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
