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

	Logger = log.New(os.Stdout, "[http] ", log.LstdFlags)

	client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: &httpclient.LoggingRoundTripper{Transport: http.DefaultTransport},
	}

	CACHE_DURATION_MINUTES = 30

	cacheInstance = cache.NewCache()
)

func randomInt(max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		Logger.Fatalf("Failed to generate random number: %v", err)
	}
	return int(nBig.Int64())
}

func LatitudeLongitude(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		handleGetRequest(w)
	case "POST":
		handlePostRequest(w, r)
	default:
		http.Error(w, `{"status": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func handleGetRequest(w http.ResponseWriter) {
	fmt.Fprintf(w, `{"status": "healthy"}`)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	var c geo.Coordinates

	if err := decodeRequestBody(r, &c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCoordinates(c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("%s,%s", c.Lat, c.Long)
	if cachedResponse, found := cacheInstance.Get(cacheKey); found {
		handleCacheHit(w, cachedResponse)
		return
	}

	handleCacheMiss(w, c, cacheKey)
}

func decodeRequestBody(r *http.Request, c *geo.Coordinates) error {
	return json.NewDecoder(r.Body).Decode(c)
}

func validateCoordinates(c geo.Coordinates) error {
	if c.Lat == "" || c.Long == "" {
		return fmt.Errorf("Lat and/or Long positions error - not set")
	}
	return nil
}

func handleCacheHit(w http.ResponseWriter, cachedResponse interface{}) {
	w.Header().Set("X-Cache-Status", "HIT")
	json.NewEncoder(w).Encode(cachedResponse)
}

func handleCacheMiss(w http.ResponseWriter, c geo.Coordinates, cacheKey string) {
	w.Header().Set("X-Cache-Status", "MISS")

	url := fmt.Sprintf(endpoint[randomInt(len(endpoint))], c.Lat, c.Long)
	resp, err := client.Get(url)
	defer resp.Body.Close()

	if err != nil {
		http.Error(w, fmt.Sprintf("HTTP request error: %s", err), http.StatusInternalServerError)
		return
	}

	if err := decodeResponseBody(resp, &location); err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %s", err), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("External API error: %s", location), resp.StatusCode)
		return
	}

	cacheInstance.Set(cacheKey, location, time.Duration(CACHE_DURATION_MINUTES)*time.Minute)
	renderTemplate(w, location)
}

func decodeResponseBody(resp *http.Response, location *interface{}) error {
	return json.NewDecoder(resp.Body).Decode(location)
}

func renderTemplate(w http.ResponseWriter, location interface{}) {
	tmpl := template.Must(template.New("safeTemplate").Parse("{{.}}"))
	if err := tmpl.Execute(w, json.NewEncoder(w).Encode(location)); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}
