FROM golang:1.11-alpine AS builder
LABEL maintainer="Said Sef <saidsef@gmail.com>"

ENV OPEN_FAAS 0.9.8

WORKDIR /app
COPY geocode.go /app/
RUN apk add --no-cache curl git && \
    go get github.com/joho/godotenv/autoload && \
    go build geocode.go && \
    curl -sL https://github.com/openfaas/faas/releases/download/${OPEN_FAAS}/fwatchdog > /usr/bin/fwatchdog && \
    chmod ag+rwx /app/geocode /usr/bin/fwatchdog

###############################################################################

FROM alpine:3.8
LABEL maintainer="Said Sef <saidsef@gmail.com>"

COPY --from=builder /app/geocode /usr/bin/
COPY --from=builder /usr/bin/fwatchdog /usr/bin/
ENV fprocess="/usr/bin/geocode"
HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["fwatchdog"]
