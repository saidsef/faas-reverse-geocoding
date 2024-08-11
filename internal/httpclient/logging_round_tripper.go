// Package httpclient provides custom HTTP client utilities.
package httpclient

import (
	"net/http"
	"time"

	"github.com/saidsef/faas-reverse-geocoding/internal/utils"
)

// LoggingRoundTripper is a custom RoundTripper that logs the details of each HTTP request and response.
type LoggingRoundTripper struct {
	Transport http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and logs the request and response details.
// It logs the HTTP method, URL, protocol, response status code, and the duration of the request.
// If an error occurs during the request, it logs the error and returns it.
func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	resp, err := lrt.Transport.RoundTrip(req)
	duration := time.Since(startTime)

	if err != nil {
		if utils.Verbose {
			utils.Logger.Errorf("HTTP request error: %s", err)
		}
		return resp, err
	}

	utils.Logger.Infof("Request: %s %s %s, Response: %d, Duration: %s, remote: %s", req.Method, req.URL, req.Proto, resp.StatusCode, duration, req.RemoteAddr)
	return resp, nil
}
