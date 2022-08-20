package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	port string
)

// Coordinates set Lat and Long
type Coordinates struct {
	Lat  string `json:"lat"`
	Long string `json:"lon"`
}

func latitudeLongitude(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var c Coordinates

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
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		client := &http.Client{Transport: tr}
		url := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1&lat=%s&lon=%s", c.Lat, c.Long)
		resp, err := client.Get(url)
		if err != nil {
			fmt.Fprintf(w, "http get error %s", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(w, "body error %s", err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprintf(w, "%s", bodyString)
		}
	default:
		fmt.Fprintf(w, `{"status": "method not allowed"}`)
	}
}

func main() {
	flag.StringVar(&port, "port", "8080", "Listening PORT")
	flag.Parse()

	r := http.NewServeMux()

	r.HandleFunc("/", latitudeLongitude)
	r.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":"+port, r))
}
