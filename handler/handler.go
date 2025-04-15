package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const PAGELIMIT = 5

type response struct {
	Status       string      `json:"status"`
	CurrentPage  int         `json:"currentPage,omitempty"`
	Limit        int         `json:"limit,omitempty"`
	TotalPages   int         `json:"totalPages,omitempty"`
	TotalRecords int         `json:"totalRecords,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	ErrMessage   string      `json:"errmessage,omitempty"`
}

// Get data from public API
func getFromExtAPI(url string) ([]byte, error) {

	// http client to make request
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for HTTP status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Standard function to handle response
func apiResponse(w http.ResponseWriter, _ *http.Request, status int, currentPage int, limit int, totalPages int, totalRecords int, data interface{}, errorMessage string) {

	var response response = response{
		Status:       http.StatusText(status),
		CurrentPage:  currentPage,
		Limit:        limit,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
		Data:         data,
		ErrMessage:   errorMessage,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	log.Printf("Responded with status %d", status)
	json.NewEncoder(w).Encode(response)
}

// Home Handler
func Home(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	apiResponse(w, nil, 200, 0, 0, 0, 0, "Server is functioning", "")
}

func ListAllObjects(w http.ResponseWriter, r *http.Request) {

	var err error

	requestedPage := r.URL.Query().Get("page")
	requestedLimit := r.URL.Query().Get("limit")

	type objectList struct {
		ID   string      `json:"id"`
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Get Object list from public API
	bodyData, err := getFromExtAPI("https://api.restful-api.dev/objects")
	if err != nil {
		panic(err)
	}
	var responseData []objectList
	err = json.NewDecoder(strings.NewReader(string(bodyData))).Decode(&responseData)
	if err != nil {
		panic(err)
	}

	// Paginate the response
	currentPage, err := strconv.Atoi(requestedPage)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	limit, err := strconv.Atoi(requestedLimit)
	if err != nil || limit < 1 {
		limit = 1
	} else if limit > PAGELIMIT {
		limit = PAGELIMIT
	}

	var offset = ((currentPage - 1) * PAGELIMIT)
	var end = offset + limit
	if end > len(responseData) {
		end = len(responseData)
	}

	var totalPages = int(math.Ceil(float64(len(responseData)) / float64(PAGELIMIT)))
	var totalRecords = len(responseData)

	responseData = responseData[offset:end]

	apiResponse(w, nil, http.StatusOK, currentPage, limit, totalPages, totalRecords, responseData, "")

}

func ListSingleObject(w http.ResponseWriter, r *http.Request) {

	var err error
	type object struct {
		ID   string      `json:"id"`
		Name string      `json:"name"`
		Data interface{} `json:"data"`
	}
	var responseData object

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	id := parts[(len(parts) - 1)]

	// Get object of particular id
	bodyData, err := getFromExtAPI("https://api.restful-api.dev/objects/" + id)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(strings.NewReader(string(bodyData))).Decode(&responseData)
	if err != nil {
		panic(err)
	}

	apiResponse(w, nil, http.StatusOK, 0, 0, 0, 0, responseData, "")

}
