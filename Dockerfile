FROM golang:1.9-alpine AS builder
MAINTAINER Said Sef <saidsef@gmail.com>

WORKDIR /app
COPY geocode.go /app/
RUN apk add --no-cache curl && \
    go build geocode.go && \
    curl -sL https://github.com/openfaas/faas/releases/download/0.7.8/fwatchdog > /usr/bin/fwatchdog && \
    chmod ag+rwx /app/geocode /usr/bin/fwatchdog

###############################################################################

FROM alpine:3.7
MAINTAINER Said Sef <saidsef@gmail.com>

COPY --from=builder /app/geocode /usr/bin/
COPY --from=builder /usr/bin/fwatchdog /usr/bin/
ENV fprocess="/usr/bin/geocode"
HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["fwatchdog"]
