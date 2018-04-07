# Golang Reverse GeoCoding

Get address from GPS coordinates

## Prerequisite

You must have Google Maps API Key

## Deployment

From GitHub:

```shell
git clone https://github.com/saidsef/faas-reverse-geocoding.git
cd faas-reverse-geocoding/
faas-cli deploy -f ./faas-reverse-geocoding.yml
```

```shell
faas-cli deploy -f https://raw.githubusercontent.com/saidsef/faas-reverse-geocoding/master/faas-reverse-geocoding.yml
```

Take it for a test drive:

```python
from requests import post
data='{"lat":"41.40338","lng":"2.17403"}'
r = post('http://localhost:8080/function/reverse-geocoding', data=data)
```
