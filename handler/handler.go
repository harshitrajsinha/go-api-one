package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// Standard function to handle response
func apiResponse(w http.ResponseWriter, _ *http.Request, status int, data interface{}) {

	var response Response = Response{
		Status: http.StatusText(status),
		Data:   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	log.Printf("Responded with status %d", status)
	json.NewEncoder(w).Encode(response)
}

func Home(w http.ResponseWriter, _ *http.Request) {
	apiResponse(w, nil, 200, "Server is functioning")
}
