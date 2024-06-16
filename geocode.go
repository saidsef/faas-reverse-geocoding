package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/saidsef/faas-reverse-geocoding/internal/handlers"
)

var (
	port    int
	verbose bool
	logger  = log.New(os.Stdout, "[http] ", log.LstdFlags)
)

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s %s %d %s %s", r.RemoteAddr, r.Method, r.URL, r.ContentLength, r.Host, r.Proto)
		next.ServeHTTP(w, r)
	}
}

func main() {
	flag.IntVar(&port, "port", 8080, "HTTP listening PORT")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	flag.Parse()

	r := http.NewServeMux()
	r.HandleFunc("/", loggingMiddleware(handlers.LatitudeLongitude))
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
