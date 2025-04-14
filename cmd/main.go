package main

import (
	"log"
	"net/http"
	"os"

	"github.com/harshitrajsinha/go-api-one/handler"
	"github.com/harshitrajsinha/go-api-one/middleware"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load()
}

func main() {

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8000"
	}

	// Routes handlers to accept HTTP requests
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/api/v1/objects", handler.ListAllObjects)
	http.HandleFunc("/api/v1/object/{id}", handler.ListSingleObject)

	// route middleware
	router := middleware.ErrorRecoveryMiddlware(middleware.LoggingMiddleware(http.DefaultServeMux))

	log.Printf("Server is listening on PORT %s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
