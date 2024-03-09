package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLatitudeLongitudeGET tests the GET request handling of the latitudeLongitude function.
// It checks if the function returns a "healthy" status as expected.
func TestLatitudeLongitudeGET(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(latitudeLongitude)

	handler.ServeHTTP(rr, req)

	expected := `{"status": "healthy"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestLatitudeLongitudePOSTValid tests the POST request handling with valid coordinates.
// It checks if the function can successfully decode the JSON payload and make an external API call.
func TestLatitudeLongitudePOSTValid(t *testing.T) {
	// Mock coordinates
	coordinates := Coordinates{
		Lat:  "51.5074",
		Long: "0.1278",
	}
	jsonPayload, err := json.Marshal(coordinates)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(latitudeLongitude)

	handler.ServeHTTP(rr, req)

	// Since the actual external API call is made, we cannot predict the exact response body.
	// However, we can check if the response status code is OK (200), indicating a successful API call.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestLatitudeLongitudePOSTInvalid tests the POST request handling with invalid (empty) coordinates.
// It checks if the function correctly handles the error case when coordinates are not provided.
func TestLatitudeLongitudePOSTInvalid(t *testing.T) {
	// Empty coordinates
	jsonPayload, err := json.Marshal(Coordinates{})
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(latitudeLongitude)

	handler.ServeHTTP(rr, req)

	expected := "Lat and/or Lon positions error - not set"
	body, _ := io.ReadAll(rr.Body)
	if !bytes.Contains(body, []byte(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(body), expected)
	}
}

// TestLatitudeLongitudeUnsupportedMethod tests the handling of unsupported HTTP methods.
// It ensures that the function returns a "method not allowed" status for methods other than GET and POST.
func TestLatitudeLongitudeUnsupportedMethod(t *testing.T) {
	req, err := http.NewRequest("PUT", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(latitudeLongitude)

	handler.ServeHTTP(rr, req)

	expected := `{"status": "method not allowed"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
