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

	// Routes handlers to accept HTTP requests
	http.HandleFunc("/", handler.Home)

	// route middleware
	router := middleware.ErrorRecoveryMiddlware(middleware.LoggingMiddleware(http.DefaultServeMux))

	log.Printf("Server is listening on PORT %s", serverPort)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
