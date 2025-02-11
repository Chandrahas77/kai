package main

import (
	"fmt"
	"kai-sec/internal/handlers"
	"log"
	"net/http"
)

func main() {
	router := handlers.NewRouter()

	// Start the server
	port := "8080"
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
