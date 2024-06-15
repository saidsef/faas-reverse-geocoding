// Package main provides a simple HTTP server that interacts with the Nominatim API
// to reverse geocode latitude and longitude coordinates. It also serves Prometheus
// metrics at the /metrics endpoint.
package main

import (
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// endpoint specifies the URLs template for the Nominatim reverse geocoding API.
	endpoint = []string{
		"https://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1&lat=%s&lon=%s",
		"https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%s&longitude=%s&localityLanguage=en",
	}

	// port defines the port on which the server listens. It can be set via command-line flag.
	port int

	// verbose defines the verbosity of logs. It can be set via command-line flag.
	verbose bool

	// location response json interface
	location interface{}

	// logger provides a logging instance prefixed with "[http]" and standard flags.
	logger = log.New(os.Stdout, "[http] ", log.LstdFlags)

	// client is an HTTP client configured with a timeout to use for external API requests.
	client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: &loggingRoundTripper{http.DefaultTransport},
	}
)

// loggingRoundTripper is a custom RoundTripper that logs the details of each HTTP request and response.
type loggingRoundTripper struct {
	transport http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and logs the request and response details.
func (lrt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	startTime := time.Now()
	resp, err := lrt.transport.RoundTrip(req)
	duration := time.Since(startTime)

	if err != nil {
		logger.Printf("HTTP request error: %s", err)
		return nil, err
	}

	logger.Printf("Request: %s %s %s, Response: %d, Duration: %s", req.Method, req.URL, req.Proto, resp.StatusCode, duration)
	return resp, nil
}

// loggingMiddleware is a middleware function that logs incoming HTTP requests.
// It logs the remote address, HTTP method, and the request URL.
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s %d %s %s", r.RemoteAddr, r.Method, r.URL, r.ContentLength, r.Host, r.Proto)
		next.ServeHTTP(w, r)
	}
}

// Coordinates defines a structure for geographical coordinates with latitude and longitude.
// It is designed to be used in applications that require geographical locations to be represented
// in a structured format. The latitude (Lat) and longitude (Long) are stored as strings to accommodate
// various formats, but they typically represent decimal degrees.
type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"lon"`
}

// latitudeLongitude is an HTTP handler function that processes incoming requests.
// It supports GET and POST methods. GET requests respond with a simple health check message.
// POST requests expect a JSON body with "lat" and "lon" fields, and return reverse geocoded
// address information using the Nominatim API.
func latitudeLongitude(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, `{"status": "healthy"}`)
	case "POST":
		var c Coordinates

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if c.Lat == "" || c.Long == "" {
			http.Error(w, "Lat and/or Long positions error - not set", http.StatusBadRequest)
			return
		}

		url := fmt.Sprintf(endpoint[randomInt(len(endpoint))], c.Lat, c.Long)
		resp, err := client.Get(url)
		defer resp.Body.Close()

		if err != nil {
			http.Error(w, fmt.Sprintf("HTTP request error: %s", err), http.StatusInternalServerError)
			return
		}

		if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
			http.Error(w, fmt.Sprintf("Error reading response body: %s", err), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("External API error: %s", location), resp.StatusCode)
			return
		}

		// Define a template that safely escapes data.
		tmpl := template.Must(template.New("safeTemplate").Parse("{{.}}"))
		if err := tmpl.Execute(w, json.NewEncoder(w).Encode(location)); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, `{"status": "method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
}

// randomInt generates a cryptographically secure random integer between 0 and max-1.
func randomInt(max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		logger.Fatalf("Failed to generate random number: %v", err)
	}
	return int(nBig.Int64())
}

// main initializes the server, setting up routes and starting the server on the specified port.
// It listens on the root path for reverse geocoding requests and on /metrics for Prometheus metrics.
func main() {
	flag.IntVar(&port, "port", 8080, "HTTP listening PORT")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	r := http.NewServeMux()
	r.HandleFunc("/", loggingMiddleware(latitudeLongitude))
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ErrorLog:          logger,
		Handler:           r,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	logger.Printf("Server is running on port %d and address %s", port, srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
