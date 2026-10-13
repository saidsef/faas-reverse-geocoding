# Installation

## Requirements

- Kubernetes 1.27 or later for the Helm chart, as `kubeVersion` in [`Chart.yaml`](https://github.com/saidsef/faas-reverse-geocoding/blob/main/charts/reverse-geocoding/Chart.yaml) sets
- Helm 3 or later for the chart, or `kubectl` with kustomize support for the manifests
- Nodes that run `linux/amd64`, the only platform the image supports
- Outbound HTTPS from the pods to `nominatim.openstreetmap.org` and `api.bigdatacloud.net`
- Go at the version `go.mod` sets, to build from source

## Helm chart

The Charts workflow publishes the chart to a Helm repository on GitHub Pages. Add the repository and install the chart:

```sh
helm repo add geocode https://saidsef.github.io/faas-reverse-geocoding
helm repo update
helm upgrade --install geocode geocode/reverse-geocoding --namespace geocode --create-namespace
```

The chart names its objects `geocode` when the release name contains `geocode`, and `<release>-geocode` otherwise. It creates these objects:

| Kind | Created when | Purpose |
|------|--------------|---------|
| Deployment | Always | Runs `replicaCount` pods, one by default, or leaves the count to the HorizontalPodAutoscaler |
| Service | Always | Exposes the pods on `service.port`, 8080 by default, as type `ClusterIP` by default |
| Ingress | `ingress.enabled` is `true` | Routes each host in `ingress.hosts` to the Service |
| HorizontalPodAutoscaler | `autoscaling.enabled` is `true` | Scales the Deployment on CPU and memory use |
| ServiceAccount | `serviceAccount.create` is `true` | Gives the pods their own identity |
| Pod `<name>-test-connection` | `helm test` runs | Requests the Service with `wget` |

[Helm values](./configuration.md#helm-values) covers the values that change what the chart deploys.

The Service has no external address by default. Forward a local port to it:

```sh
kubectl port-forward -n geocode svc/geocode 8080:8080
```

Send a request from a second terminal:

```sh
curl -d '{"lat":"41.40338","lon":"2.17403"}' http://localhost:8080/
```

## Kustomize

The manifests live under [`deployment/`](https://github.com/saidsef/faas-reverse-geocoding/tree/main/deployment). No kustomization in it sets a namespace, so pass one to `kubectl`:

```sh
git clone https://github.com/saidsef/faas-reverse-geocoding.git
cd faas-reverse-geocoding
kubectl create namespace geocode
kubectl apply -n geocode -k deployment/
```

The base in `deployment/base/` creates these objects, both named `geocode`:

| Kind | Purpose |
|------|---------|
| Deployment | One replica of `ghcr.io/saidsef/faas-reverse-geocoding:latest`, with `imagePullPolicy: Always`, so the kubelet checks the registry each time it starts the container |
| Service | `ClusterIP` on port 8080, with the port named `tcp-term` and `ClientIP` session affinity for 60 seconds |

Each of these directories builds on the base:

| Directory | Adds |
|-----------|------|
| `deployment/` | The `app.kubernetes.io/part-of: geocode` annotation |
| `deployment/gateway-api/` | An HTTPRoute that sends paths under `/geocode` to the Service, attached to the `https` listener of the Gateway `external-gw` in the namespace `gateway-ingress` |
| `deployment/ingress-nginx/` | An Ingress for the host `geocode`, with CORS on and a configuration snippet for `/metrics`, described in [Metrics](./metrics.md#scraping) |

The `gateway-api` overlay does not create the Gateway. Edit `parentRefs` in [`route.yml`](https://github.com/saidsef/faas-reverse-geocoding/blob/main/deployment/gateway-api/route.yml) to name a Gateway in your cluster before you apply it. The service answers on any path, so the route needs no rewrite.

## Pod security

The chart and the base run the container with the same restrictions:

- User and group 65534, with `runAsNonRoot` set
- A read-only root filesystem
- No Linux capabilities and no privilege escalation
- The `RuntimeDefault` seccomp profile
- No service account token and no service environment variables in the pod

The liveness probe checks that the `faas-reverse-geocoding` process runs, and the readiness probe opens a TCP connection to port 8080.

## Container image

CI publishes the image to `ghcr.io/saidsef/faas-reverse-geocoding` for `linux/amd64` only. A pod on an `arm64` node fails to pull it with `no match for platform in manifest`.

| Build | Tags |
|-------|------|
| A push to `main` | `latest`, and `vYYYY.MM` for the year and month of the build |
| Any other branch | The branch name, with each `.` and `_` turned into `-`, cut to 50 characters |

The image runs the binary as the user `nobody` on Alpine Linux. Its `HEALTHCHECK` requests `/` on port 8080 with `wget` every 60 seconds.

## Running locally

```sh
docker run --rm -p 8080:8080 ghcr.io/saidsef/faas-reverse-geocoding:latest
```

On an `arm64` machine, add `--platform linux/amd64`, which runs the image under emulation. To run from a clone of the repository instead:

```sh
go run . -verbose
```

## Building from source

```sh
go build -o faas-reverse-geocoding .
docker build -t faas-reverse-geocoding .
```

`make build` writes the same binary to `bin/faas-reverse-geocoding`.

## Log output

```sh
kubectl logs -n geocode deploy/geocode
```

The service writes one JSON object per line to standard output. The start-up line records the port:

```json
{"level":"info","msg":"Server is running on port 8080 and address :8080","time":"2026-10-13T10:06:01Z"}
```

Each request on any path except `/metrics` writes a line with the client address, the method, the URL path and query, the body length in bytes, the `Host` header and the protocol:

```json
{"level":"info","msg":"192.168.65.1:24366 POST / 22 localhost:8080 HTTP/1.1","time":"2026-10-13T10:06:02Z"}
```

Each cache miss writes a line for the upstream request with this message, or an error line naming the URL when the request fails:

```text
Request: GET <upstream URL> HTTP/1.1, Response: <status>, Duration: <time>, remote:
```

The upstream URL carries the coordinates, so the log holds every location that missed the cache. [`-verbose`](./configuration.md#command-line-flags) adds debug lines for each cache lookup and error lines for rejected requests.
