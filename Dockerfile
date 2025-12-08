FROM docker.io/golang:1.24-alpine3.20 AS builder
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"

WORKDIR /app
COPY . /app/
RUN apk add --no-cache curl git && \
    go get -v && \
    go build -v -o /app/ ./...

FROM docker.io/alpine:3.23
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"
ENV PORT=${PORT:-8080}

USER nobody

COPY --from=builder /app/faas-reverse-geocoding /usr/bin/

HEALTHCHECK --interval=60s --timeout=10s CMD curl --fail 'http://localhost:${PORT}/' || exit 1

EXPOSE ${PORT}

CMD ["/usr/bin/faas-reverse-geocoding"]
