FROM golang:1.9-alpine AS builder
MAINTAINER Said Sef <saidsef@gmail.com>

WORKDIR /app
ADD https://github.com/openfaas/faas/releases/download/0.7.7/fwatchdog /usr/bin
COPY function.go /app/
RUN go build function.go && \
    chmod ag+rwx /app/function /usr/bin/fwatchdog

###############################################################################

FROM alpine
MAINTAINER Said Sef <saidsef@gmail.com>

COPY --from=builder /app/function /usr/bin/
COPY --from=builder /usr/bin/fwatchdog /usr/bin/

WORKDIR /app

ENV fprocess="/usr/bin/function"

HEALTHCHECK --interval=1s CMD [ -e /tmp/.lock ] || exit 1

CMD ["/usr/bin/fwatchdog"]
