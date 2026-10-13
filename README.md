# Reverse Geocoding

[![CI](https://github.com/saidsef/faas-reverse-geocoding/actions/workflows/ci.yml/badge.svg)](https://github.com/saidsef/faas-reverse-geocoding/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/saidsef/faas-reverse-geocoding)](https://goreportcard.com/report/github.com/saidsef/faas-reverse-geocoding)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/saidsef/faas-reverse-geocoding)
[![GoDoc](https://godoc.org/github.com/saidsef/faas-reverse-geocoding?status.svg)](https://pkg.go.dev/github.com/saidsef/faas-reverse-geocoding?tab=doc)
![GitHub release(latest by date)](https://img.shields.io/github/v/release/saidsef/faas-reverse-geocoding)
![Commits](https://img.shields.io/github/commits-since/saidsef/faas-reverse-geocoding/latest.svg)
![GitHub](https://img.shields.io/github/license/saidsef/faas-reverse-geocoding)

Reverse geocoding turns a latitude and longitude into an address. This service takes the coordinates as JSON, asks OpenStreetMap's Nominatim or BigDataCloud for the address at that point, and returns the provider's answer. It keeps each answer in memory, for 60 minutes by default, so repeat lookups of the same place skip the provider.

- **Lookup** - takes `lat` and `lon` in a JSON body and returns the address record for that point
- **Cache** - answers repeat lookups that round to the same 3 decimal places from memory
- **Health** - answers `GET` on any path with `{"status": "healthy"}`
- **Metrics** - counts upstream requests per provider on a Prometheus endpoint

## Quick start

```sh
helm repo add geocode https://saidsef.github.io/faas-reverse-geocoding
helm repo update
helm upgrade --install geocode geocode/reverse-geocoding --namespace geocode --create-namespace
kubectl port-forward -n geocode svc/geocode 8080:8080
```

From a second terminal:

```sh
curl -d '{"lat":"41.40338","lon":"2.17403"}' http://localhost:8080/
```

## Documentation

[Read the Docs](https://faas-reverse-geocoding.readthedocs.io/en/latest/) hosts the same pages.

| Page | Contents |
|------|----------|
| [Installation](https://github.com/saidsef/faas-reverse-geocoding/blob/main/docs/installation.md) | Requirements, the Helm chart, the Kustomize overlays, the container image, running locally, building from source, log output |
| [Configuration](https://github.com/saidsef/faas-reverse-geocoding/blob/main/docs/configuration.md) | Command-line flags, setting flags in Kubernetes, Helm values |
| [API](https://github.com/saidsef/faas-reverse-geocoding/blob/main/docs/api.md) | Endpoints, the request body, responses, errors, coordinate precision |
| [Architecture](https://github.com/saidsef/faas-reverse-geocoding/blob/main/docs/architecture.md) | What one request does, the cache, the upstream providers, timeouts, logging |
| [Metrics](https://github.com/saidsef/faas-reverse-geocoding/blob/main/docs/metrics.md) | Prometheus endpoint, metric names and labels, scrape setup |

## Requirements

The Helm chart needs Kubernetes 1.27 or later, and the image runs on `linux/amd64` nodes only. Building from source needs the Go version that `go.mod` sets.

## Source

Our latest and greatest source of *Reverse Geocoding* can be found on [GitHub](https://github.com/saidsef/faas-reverse-geocoding). [Fork us](https://github.com/saidsef/faas-reverse-geocoding/fork)!

## Contributing

We would :heart: you to contribute by making a [pull request](https://github.com/saidsef/faas-reverse-geocoding/pulls).

Please read the official [Contribution Guide](https://github.com/saidsef/faas-reverse-geocoding/blob/main/CONTRIBUTING.md) for more information on how you can contribute.
