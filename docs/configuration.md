# Configuration

The binary takes three command-line flags and reads no environment variables. The Helm chart has values of its own, which change what the chart deploys but not how the binary behaves.

## Command-line flags

| Flag | Values | Default | Description |
|------|--------|---------|-------------|
| `-port` | A TCP port number | `8080` | The port for every endpoint, `/metrics` included |
| `-cache` | A whole number of minutes | `60` | How long a successful lookup stays in the cache. With `0` or less, every request goes to the upstream provider |
| `-verbose` | `true` or `false`, and `-verbose` on its own means `true` | `false` | Sets the log level to `debug`, which adds a line for each cache lookup and the error behind each rejected request |

Each pod keeps its own cache. [Cache](./architecture.md#cache) describes how entries expire.

## Setting flags in Kubernetes

Neither the Kustomize base nor the Helm chart passes any arguments, so the binary runs with the defaults above. The image sets `CMD` and no `ENTRYPOINT`, so a container that sets `args` alone fails to start with `exec: "-cache": executable file not found in $PATH`. Set `command` as well:

```yaml
containers:
  - name: geocode
    command: ["/usr/bin/faas-reverse-geocoding"]
    args: ["-cache", "30", "-verbose"]
```

The chart has no value for `command` or `args`, so a chart install always runs with the defaults.

The manifests, the chart and the image `HEALTHCHECK` all expect port 8080. A different `-port` needs the matching change to `containerPort`, the probes and the Service, and the chart template fixes `containerPort` at 8080.

## Helm values

[`values.yaml`](https://github.com/saidsef/faas-reverse-geocoding/blob/main/charts/reverse-geocoding/values.yaml) holds every value and its default. These values behave in ways the file does not show:

| Value | Behaviour |
|-------|-----------|
| `image.tag` | Defaults to `latest` with `image.pullPolicy` set to `IfNotPresent`, so a node keeps the first `latest` it pulled. A `vYYYY.MM` tag from [Container image](./installation.md#container-image) pins the version. An empty tag falls back to the chart's `appVersion`, which CI never publishes as an image tag |
| `replicaCount` | Has no effect while `autoscaling.enabled` is `true` |
| `service.port` | Sets the Service port only. The Service targets the container port by its name, `http`, and the container port stays 8080 |
| `ingress.hosts` | The template reads `host` from each entry and ignores `paths`. Every host routes `/` with `pathType: Prefix` |
| `ingress.className` | Sets `spec.ingressClassName` on the Ingress |
| `autoscaling.enabled` | Creates an `autoscaling/v2` HorizontalPodAutoscaler between `minReplicas` and `maxReplicas`, with targets for average CPU and memory use as a percentage of the pod's requests |
| `serviceAccount.create` | Defaults to `false`, and the pods run as the namespace's `default` service account. With `true`, the chart creates a ServiceAccount named `serviceAccount.name` and runs the pods as it. Either way the pods get no service account token |
| `podAnnotations` | Defaults to `prometheus.io/scrape: "true"` and `filter.by.port.name: "true"`, described in [Metrics](./metrics.md#scraping) |
| `livenessProbe`, `readinessProbe` | Replace the probes whole |
| `nodeSelector` | Empty by default. The image runs on `linux/amd64` only, so a cluster with `arm64` nodes needs `kubernetes.io/arch: amd64` here |
