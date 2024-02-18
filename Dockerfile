FROM docker.io/golang:1.22-alpine3.19 AS builder
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"

ENV PORT ${PORT:-8080}

WORKDIR /app
COPY geocode.go go.mod go.sum /app/
RUN apk add --no-cache curl git && \
    go get github.com/prometheus/client_golang/prometheus/promhttp && \
    go build geocode.go

FROM docker.io/alpine:3.19
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"

USER nobody

COPY --from=builder /app/geocode /usr/bin/

HEALTHCHECK --interval=60s --timeout=10s CMD curl --fail 'http://localhost:${PORT}/' || exit 1

EXPOSE ${PORT}

CMD ["/usr/bin/geocode"]
