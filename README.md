# Golang Reverse GeoCoding

Reverse geocoding is used to find places or addresses near a latitude, longitude pair. Picture a map showing building outlines but no labels, then clicking on a building and being shown the name of the business. That is reverse geocoding.

This repository holds a small web service that performs reverse geocoding to determine whether specified geo  coordinates has an adress. If it is then the response will contain attributes associated with the matched adress, municipality, county, country, etc.

## Deployment

### Helm Chart

```shell
helm repo add geocode https://saidsef.github.io/faas-reverse-geocoding
helm repo update 
helm upgrade --install geocode geocode/reverse-geocoding --namespace geocode --create-namespace
```

> *NOTE:* API can be accessed via port-forward `Service` or via Enabling `Ingress`

### Kustomization

```shell
kubectl apply -k deployment/
```

> *NOTE:* API can be accessed via port-forward `Service` or via updating `Ingress`

Take it for a test drive:

```shell
curl -d '{"lat":"41.40338","lon":"2.17403"}' http://localhost:8080/
```

```python
from requests import post
data='{"lat":"41.40338","lon":"2.17403"}'
r = post('http://localhost:8080/', data=data)
print(r.text)
```

## Source

Our latest and greatest source of Jenkins can be found on [GitHub]. Fork us!

## Contributing

We would :heart: you to contribute by making a [pull request](https://github.com/saidsef/faas-reverse-geocoding/pulls).

Please read the official [Contribution Guide](./CONTRIBUTING.md) for more information on how you can contribute.
