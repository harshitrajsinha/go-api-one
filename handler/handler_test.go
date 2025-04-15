package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Unit test for Home route
func TestHome(t *testing.T) {

	var actualResponse map[string]string
	var err error

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	Home(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", res.StatusCode)
	}

	expectedBody := map[string]string{
		"status": "OK",
		"data":   "Server is functioning",
	}

	// Read resposne body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&actualResponse)

	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// check if expected and acutal output are same
	if !cmp.Equal(actualResponse, expectedBody) {
		t.Errorf("Expected response %v, got %v", expectedBody, actualResponse)
	}

}

// Unit test for list of objects with default page and limit
func TestListAllObjects(t *testing.T) {

	var err error

	// Type cast integer to float
	var expectedResponse = map[string]interface{}{
		"status":       "OK",
		"currentPage":  1.0,
		"limit":        1.0,
		"totalPages":   3.0,
		"totalRecords": 13.0,
		"data": []interface{}{
			map[string]interface{}{
				"id":   "1",
				"name": "Google Pixel 6 Pro",
				"data": map[string]interface{}{
					"color":    "Cloudy White",
					"capacity": "128 GB",
				},
			},
		},
	}

	var actualResponse map[string]interface{}

	req := httptest.NewRequest(http.MethodGet, "/objects", nil)
	w := httptest.NewRecorder()

	ListAllObjects(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error while reading response body %v", err)
	}

	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&actualResponse)
	if err != nil {
		t.Errorf("Failed to convert response data to valid structure %v", err)
	}

	if diff := cmp.Diff(expectedResponse, actualResponse); diff != "" {
		t.Errorf("Mismatch (-expected +actual):\n%s", diff)
	}

	// check if expected and acutal output are same
	if !cmp.Equal(actualResponse, expectedResponse) {
		t.Errorf("Expected response %v, got %v", expectedResponse, actualResponse)
	}

}

// Test for single object with ID = 7
func TestListSingleObject(t *testing.T) {

	var err error
	var actualResponse map[string]interface{}
	var expectedResponse = map[string]interface{}{
		"status": "OK",
		"data": map[string]interface{}{
			"id":   "7",
			"name": "Apple MacBook Pro 16",
			"data": map[string]interface{}{
				"year":           2019.0,
				"price":          1849.99,
				"CPU model":      "Intel Core i9",
				"Hard disk size": "1 TB",
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/object/7", nil)
	w := httptest.NewRecorder()

	ListSingleObject(w, req)

	res := w.Result()
	defer res.Body.Close()

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error reading response body %v", body)
	}

	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&actualResponse)
	if err != nil {
		t.Errorf("Failed to convert response data to valid structure %v", err)
	}

	if diff := cmp.Diff(expectedResponse, actualResponse); diff != "" {
		t.Errorf("Mismatch (-expected +actual):\n%s", diff)
	}

	// check if expected and acutal output are same
	if !cmp.Equal(actualResponse, expectedResponse) {
		t.Errorf("Expected response %v, got %v", expectedResponse, actualResponse)
	}

}
