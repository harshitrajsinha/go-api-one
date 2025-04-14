package main

import (
	"encoding/json"
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
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Welcome to home route")
	})

	// route middleware
	router := middleware.ErrorRecoveryMiddlware(middleware.LoggingMiddleware(http.DefaultServeMux))

	log.Printf("Server is listening on PORT %s", serverPort)
	if err := http.ListenAndServe(":"+serverPort, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
