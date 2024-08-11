// Package metrics provides functionality for collecting and exposing
// metrics related to host requests using Prometheus.
package metrics

import "github.com/prometheus/client_golang/prometheus"

// Hostname is a counter vector that tracks the total number of requests
// to each host. It is labelled by the hostname to differentiate between
// different hosts.
var (
	Hostname = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "host_request_total",
			Help: "Total number of requests to host.",
		},
		[]string{"hostname"},
	)
)

// Init registers the Hostname counter vector with Prometheus. This function
// should be called during the initialisation phase of the application to
// ensure that the metrics are available for scraping by Prometheus.
func Init() {
	prometheus.MustRegister(Hostname)
}
