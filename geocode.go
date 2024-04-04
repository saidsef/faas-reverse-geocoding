// Package main provides a simple HTTP server that interacts with the Nominatim API
// to reverse geocode latitude and longitude coordinates. It also serves Prometheus
// metrics at the /metrics endpoint.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// endpoint specifies the URL template for the Nominatim reverse geocoding API.
const (
	endpoint = "https://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1&lat"
)

var (
	// port defines the port on which the server listens. It can be set via command-line flag.
	port int

	// logger provides a logging instance prefixed with "[http]" and standard flags.
	logger = log.New(os.Stdout, "[http] ", log.LstdFlags)

	// client is an HTTP client configured with a timeout to use for external API requests.
	client = &http.Client{
		Timeout: time.Second * 10,
	}
)

// loggingMiddleware is a middleware function that logs incoming HTTP requests.
// It logs the remote address, HTTP method, and the request URL.
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s %d %s", r.RemoteAddr, r.Method, r.URL, r.ContentLength, r.Host)
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
			http.Error(w, "Lat and/or Lon positions error - not set", http.StatusBadRequest)
			return
		}

		url := fmt.Sprintf("%s=%s&lon=%s", endpoint, c.Lat, c.Long)
		resp, err := client.Get(url)
		if err != nil {
			http.Error(w, fmt.Sprintf("HTTP request error: %s", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading response body: %s", err), http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("External API error: %s", string(bodyBytes)), resp.StatusCode)
			return
		}

		fmt.Fprintf(w, "%s", bodyBytes)
	default:
		http.Error(w, `{"status": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// main initializes the server, setting up routes and starting the server on the specified port.
// It listens on the root path for reverse geocoding requests and on /metrics for Prometheus metrics.
func main() {
	flag.IntVar(&port, "port", 8080, "HTTP Listening PORT")
	flag.Parse()

	r := http.NewServeMux()
	r.HandleFunc("/", loggingMiddleware(latitudeLongitude))
	r.Handle("/metrics", promhttp.Handler())

	logger.Printf("Server is running on port %d", port)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 15 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
