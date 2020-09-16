FROM golang:1.15-alpine AS builder
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"

ENV OPEN_FAAS 0.18.18

WORKDIR /app
COPY geocode.go /app/
RUN apk add --no-cache curl git && \
    go get github.com/joho/godotenv/autoload && \
    go build geocode.go && \
    curl -sL https://github.com/openfaas/faas/releases/download/${OPEN_FAAS}/fwatchdog > /usr/bin/fwatchdog && \
    chmod a+rwx /app/geocode /usr/bin/fwatchdog

FROM alpine:3.8
LABEL maintainer="Said Sef <saidsef@gmail.com> (saidsef.co.uk/)"

ENV fprocess="/usr/bin/geocode"

COPY --from=builder /app/geocode /usr/bin/
COPY --from=builder /usr/bin/fwatchdog /usr/bin/

USER nobody

HEALTHCHECK --interval=10s CMD [ -e /tmp/.lock ] || exit 1

CMD ["/usr/bin/fwatchdog"]
