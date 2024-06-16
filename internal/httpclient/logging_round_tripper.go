// Package httpclient provides custom HTTP client utilities.
package httpclient

import (
	"log"
	"net/http"
	"os"
	"time"
)

// logger provides a logging instance prefixed with "[http]" and standard flags.
var logger = log.New(os.Stdout, "[http] ", log.LstdFlags)

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
		logger.Printf("HTTP request error: %s", err)
		return nil, err
	}

	logger.Printf("Request: %s %s %s, Response: %d, Duration: %s", req.Method, req.URL, req.Proto, resp.StatusCode, duration)
	return resp, nil
}
