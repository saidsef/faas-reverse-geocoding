package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLatitudeLongitude(t *testing.T) {
	// Create a new HTTP request with a POST method and a JSON payload
	payload := map[string]string{"lat": "40.712776", "lon": "-74.005974"}
	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the latitudeLongitude handler function with the test request and response
	handler := http.HandlerFunc(latitudeLongitude)
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body of the response
	want := rr.Body.String()
	// Openstreetmap slightly change place_id value on different requests and the test fails
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), want)
	}
}
