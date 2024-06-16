package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/saidsef/faas-reverse-geocoding/internal/cache"
	"github.com/saidsef/faas-reverse-geocoding/internal/geo"
	"github.com/saidsef/faas-reverse-geocoding/internal/httpclient"
)

var (
	endpoint = []string{
		"https://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1&lat=%s&lon=%s",
		"https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=%s&longitude=%s&localityLanguage=en",
	}

	location interface{}

	logger = log.New(os.Stdout, "[http] ", log.LstdFlags)

	client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: &httpclient.LoggingRoundTripper{Transport: http.DefaultTransport},
	}

	CACHE_DURATION_MINUTES = 30

	cacheInstance = cache.NewCache()
)

func LatitudeLongitude(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		fmt.Fprintf(w, `{"status": "healthy"}`)
	case "POST":
		var c geo.Coordinates

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if c.Lat == "" || c.Long == "" {
			http.Error(w, "Lat and/or Long positions error - not set", http.StatusBadRequest)
			return
		}

		cacheKey := fmt.Sprintf("%s,%s", c.Lat, c.Long)
		if cachedResponse, found := cacheInstance.Get(cacheKey); found {
			w.Header().Set("X-Cache-Status", "HIT")
			json.NewEncoder(w).Encode(cachedResponse)
			return
		} else {
			w.Header().Set("X-Cache-Status", "MISS")
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

		cacheInstance.Set(cacheKey, location, time.Duration(CACHE_DURATION_MINUTES)*time.Minute)

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

func randomInt(max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		logger.Fatalf("Failed to generate random number: %v", err)
	}
	return int(nBig.Int64())
}
