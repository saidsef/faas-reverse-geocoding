# Metrics

The service serves Prometheus metrics at `/metrics` on the port [`-port`](./configuration.md#command-line-flags) sets, 8080 by default. The endpoint needs no authentication and answers any method.

```sh
kubectl port-forward -n geocode svc/geocode 8080:8080
curl http://localhost:8080/metrics
```

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `host_request_total` | counter | `hostname` | Requests the service sent to an upstream provider |
| `promhttp_metric_handler_requests_total` | counter | `code` | Scrapes of `/metrics`, by HTTP status |
| `promhttp_metric_handler_requests_in_flight` | gauge | - | Scrapes of `/metrics` in progress |
| `go_*` | various | - | Go runtime statistics from the Prometheus client library |
| `process_*` | various | - | CPU, memory and file descriptor use of the service process |

`hostname` holds `nominatim.openstreetmap.org` or `api.bigdatacloud.net`. The counter moves once for each cache miss, after the upstream request returns and before the service checks the result, so failed requests count too. Cache hits and rejected requests never reach it. Until the first cache miss, the endpoint carries no `host_request_total` series at all.

The service exports no metric for client requests, response status or latency.

## Scraping

The two ways to deploy put the same annotations on different objects:

| Deployment | Object | Annotations |
|------------|--------|-------------|
| Kustomize base | The Service | `prometheus.io/scrape: "true"`, `filter.by.port.name: "true"` |
| Helm chart | Each pod, from `podAnnotations` | `prometheus.io/scrape: "true"`, `filter.by.port.name: "true"` |

Prometheus acts on these annotations only through a scrape config that reads them. Neither sets `prometheus.io/port` or `prometheus.io/path`.

The Prometheus Operator ignores the annotations. Under the operator, this ServiceMonitor selects the Service by label:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: geocode
  namespace: geocode
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: geocode
  endpoints:
    - port: tcp-term  # http for the Helm chart
```

The Service port is named `tcp-term` in the Kustomize base and `http` in the Helm chart.

The `ingress-nginx` overlay carries a `configuration-snippet` annotation with a `location /metrics` block. ingress-nginx reads snippet annotations only when the controller sets [`allow-snippet-annotations`](https://kubernetes.github.io/ingress-nginx/user-guide/nginx-configuration/configmap/#allow-snippet-annotations) to `true`, and its default is `false`.

## Queries

```promql
# upstream requests per provider over the last hour
sum by (hostname) (increase(host_request_total[1h]))

# share of upstream requests that went to Nominatim over the last day
sum(increase(host_request_total{hostname="nominatim.openstreetmap.org"}[1d]))
  / sum(increase(host_request_total[1d]))
```
