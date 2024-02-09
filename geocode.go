package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Coordinates set Lat and Long in json payload
type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"lon"`
}

var (
	// Coordinates: struct
	c Coordinates
	// Port number: string
	port string
	// Create a new logger object
	logger = log.New(os.Stdout, "[http] ", log.LstdFlags)
	// Endpoint: string
	endpoint string = "https://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1&lat"
)

// Wrap handler function with logging middleware
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	}
}

// Lat and Long geo coordinates
func latitudeLongitude(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, `{"status": "healthy"}`)
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(c.Lat) == 0 && len(c.Long) == 0 {
			fmt.Fprintf(w, "Lat and/or Lon positions error - not set")
			return
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		}

		client := &http.Client{Transport: tr}
		url := fmt.Sprintf("%s=%s&lon=%s", endpoint, c.Lat, c.Long)

		// Log outgoing request
		logger.Printf("%s %s", r.RemoteAddr, url)

		resp, err := client.Get(url)
		if err != nil {
			logger.Printf("http get error: %s", err)
			fmt.Fprintf(w, "http get error %s", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Printf("body error: %s", err)
				fmt.Fprintf(w, "body error %s", err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprintf(w, "%s", bodyString)

			// Log incoming response
			logger.Printf("%s %d %s", r.RemoteAddr, resp.StatusCode, bodyString)
		} else {
			logger.Printf("%s %d", r.RemoteAddr, resp.StatusCode)
		}
	default:
		fmt.Fprintf(w, `{"status": "method not allowed"}`)
	}
}

// Main func server
func main() {
	flag.StringVar(&port, "port", "8080", "Listening PORT")
	flag.Parse()

	r := http.NewServeMux()

	// Wrap handler function with logging middleware
	r.HandleFunc("/", loggingMiddleware(latitudeLongitude))
	r.Handle("/metrics", promhttp.Handler())

	logger.Printf("Server started on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
