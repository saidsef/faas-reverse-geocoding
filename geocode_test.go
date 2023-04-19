package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(latitudeLongitude)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"status": "healthy"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPostRequest(t *testing.T) {
	c := Coordinates{Lat: "40.748817", Long: "-73.985428"}
	payload, _ := json.Marshal(c)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(latitudeLongitude)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if strings.TrimSpace(rr.Body.String()) == "" {
		t.Errorf("handler returned empty body")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	if _, ok := response["address"]; !ok {
		t.Errorf("handler returned unexpected response: %v", response)
	}
}
