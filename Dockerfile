FROM golang:1.19-alpine AS builder
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"

ENV PORT ${PORT:-8080}

WORKDIR /app
COPY geocode.go /app/
COPY go.mod /app/
COPY go.sum /app/
RUN apk add --no-cache curl git && \
    go get github.com/prometheus/client_golang/prometheus/promhttp && \
    go build geocode.go

FROM alpine:3.17
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"

ENV fprocess="/usr/bin/geocode"

COPY --from=builder /app/geocode /usr/bin/

USER nobody

HEALTHCHECK --interval=30s --timeout=10s CMD curl --fail 'http://localhost:${PORT}/' || exit 1

EXPOSE ${PORT}

CMD ["/usr/bin/geocode"]
