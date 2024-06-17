package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saidsef/faas-reverse-geocoding/internal/handlers"
)

// TestLatitudeLongitudeHandler tests the latitudeLongitude HTTP handler function.
func TestLatitudeLongitudeHandler(t *testing.T) {
	// Test cases to cover all critical paths.
	tests := []struct {
		name           string
		method         string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Health Check - GET Request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status": "healthy"}`,
		},
		{
			name:           "POST Request - Missing Body",
			method:         "POST",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "EOF\n",
		},
		{
			name:           "POST Request - Incomplete Data",
			method:         "POST",
			body:           `{"lat": "51.5074"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "lat and/or Long positions error - not set\n",
		},
		// Add more test cases as needed, especially for successful POST requests.
		// Note: Successful POST requests would require mocking the external API call to Nominatim.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the specified method, URL, and body.
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(tc.method, "/", bytes.NewBufferString(tc.body))
			} else {
				req = httptest.NewRequest(tc.method, "/", nil)
			}
			// Record the response.
			rr := httptest.NewRecorder()

			// Create an HTTP handler from our function.
			handler := http.HandlerFunc(handlers.LatitudeLongitude)

			// Serve the HTTP request to our handler.
			handler.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// Check the response body is what we expect.
			if rr.Body.String() != tc.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tc.expectedBody)
			}
		})
	}
}
