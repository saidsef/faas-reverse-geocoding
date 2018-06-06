#!/usr/bin/env sh

set -e pipefail

BUILD_ID=$1

if ! [ -x "$(command -v docker)" ]; then
  echo 'Unable to find docker command, please install Docker (https://www.docker.com/) and retry' >&2
  exit 1
fi

if ! [ -x "${BUILD_ID}" ]; then
  echo "Build ID is missing"
  exit 1
fi

build() {
  echo "Building image"
  docker build -t saidsef/faas-reverse-geocoding:${BUILD_ID} .
  docker tag saidsef/faas-reverse-geocoding:${BUILD_ID} saidsef/faas-reverse-geocoding:latest
}

push() {
  echo "Pushing image"
  docker push saidsef/faas-reverse-geocoding
}

main() {
  build
  push
}

main
