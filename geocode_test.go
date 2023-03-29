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
	want := `{"place_id":226470475,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. https://osm.org/copyright","osm_type":"way","osm_id":575213527,"lat":"40.71273945","lon":"-74.00593904130275","display_name":"New York City Hall, 260, Broadway, Lower Manhattan, Manhattan Community Board 1, Manhattan, New York County, City of New York, New York, 10000, United States","address":{"amenity":"New York City Hall","house_number":"260","road":"Broadway","quarter":"Lower Manhattan","neighbourhood":"Manhattan Community Board 1","suburb":"Manhattan","county":"New York County","city":"City of New York","state":"New York","ISO3166-2-lvl4":"US-NY","postcode":"10000","country":"United States","country_code":"us"},"boundingbox":["40.712445","40.7130254","-74.0064455","-74.0055687"]}`
	if rr.Body.String() != want {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), want)
	}
}
