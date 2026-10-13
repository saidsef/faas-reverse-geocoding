# API

The service listens on the port [`-port`](./configuration.md#command-line-flags) sets, 8080 by default. `/metrics` serves [Prometheus metrics](./metrics.md) to any method. Every other path goes to the geocoding handler, where the method alone decides the response:

| Method | Response |
|--------|----------|
| `POST` | The address record for the coordinates in the body, described in [Reverse geocoding](#reverse-geocoding) |
| `GET` | `{"status": "healthy"}` with status 200 |
| Any other, `HEAD` included | `{"status": "method not allowed"}` with status 405 |

## Reverse geocoding

```sh
curl -d '{"lat":"41.40338","lon":"2.17403"}' http://localhost:8080/
```

```python
from requests import post

r = post("http://localhost:8080/", data='{"lat":"41.40338","lon":"2.17403"}')
print(r.status_code, r.headers.get("X-Cache-Status"), r.text)
```

### Request body

The body is a JSON object, and the service reads it whatever the `Content-Type` header says:

| Field | Type | Accepted values |
|-------|------|-----------------|
| `lat` | String | A decimal latitude from `-90` to `90` |
| `lon` | String | A decimal longitude from `-180` to `180` |

Both fields are required and both must be JSON strings. A JSON number in either field returns 400. The service parses each string as a 32-bit float, rejects surrounding spaces, and ignores any other field in the object.

### Response

A successful lookup returns status 200 with these headers:

| Header | Value |
|--------|-------|
| `Content-Type` | `application/json` |
| `X-Cache-Status` | `MISS` when the service asked an upstream provider, `HIT` when it answered from the [cache](./architecture.md#cache) |

The body is the JSON that one of two upstream providers returned. On each cache miss the service picks a provider at random, and the two return different fields:

| Provider | Request | Response reference |
|----------|---------|--------------------|
| Nominatim | `https://nominatim.openstreetmap.org/reverse?format=json&zoom=18&addressdetails=1&lat=…&lon=…` | [Nominatim reverse geocoding](https://nominatim.org/release-docs/latest/api/Reverse/) |
| BigDataCloud | `https://api.bigdatacloud.net/data/reverse-geocode-client?latitude=…&longitude=…&localityLanguage=en` | [Reverse geocode to city API](https://www.bigdatacloud.com/free-api/free-reverse-geocode-to-city-api) |

A client cannot choose the provider. A cache hit returns whichever provider answered the first request for that location.

The service decodes the provider's JSON and encodes it again. Object keys come back in alphabetical order, and `<`, `>` and `&` inside strings come back as `<`, `>` and `&`.

### Errors

Every error response has `Content-Type: text/plain; charset=utf-8` and a body that ends in a newline. An error after a cache miss keeps `X-Cache-Status: MISS`. The service caches no error.

| Status | Body | Cause |
|--------|------|-------|
| 400 | `EOF` | The body is empty |
| 400 | `invalid character 'h' looking for beginning of value`, or another JSON decoder error | The body is not JSON |
| 400 | `json: cannot unmarshal number into Go struct field .lat of type string` | `lat` or `lon` is a JSON number |
| 400 | `strconv.ParseFloat: parsing "…": invalid syntax`, quoting the value | `lat` or `lon` is missing, which quotes an empty value, or holds something other than a decimal number |
| 400 | `Latitude and/or longitude values are out of range.` | `lat` or `lon` lies outside its range |
| 405 | `{"status": "method not allowed"}` | The method is neither `GET` nor `POST` |
| 500 | `Failed to fetch data` | The upstream request failed, or took longer than 10 seconds |
| 500 | `Failed to decode response` | The upstream body is not JSON, whatever its status |
| The upstream status | `External API error: ` followed by the upstream JSON in Go syntax, such as `map[error:rate limited]` | The upstream returned JSON with a status other than 200 |

## Coordinate precision

The cache key rounds both values to 3 decimal places. Within one pod, every request that rounds to the same pair shares one cache entry, so `41.4034,2.1736` gets the answer cached for `41.40338,2.17403`. The request to the provider carries 6 decimal places of the 32-bit value, so `41.40338` goes out as `41.403381`.

One degree of latitude spans about 111 km. A degree of longitude spans the same at the equator and less towards the poles. Each decimal place divides the step by ten:

| Decimal places | Step in latitude |
|:--------------:|:----------------:|
| 1 | 11.1 km |
| 2 | 1.11 km |
| 3 | 111 m |
| 4 | 11.1 m |
| 5 | 1.11 m |
| 6 | 11.1 cm |
