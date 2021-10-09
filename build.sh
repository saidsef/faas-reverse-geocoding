#!/usr/bin/env sh

set -e

BUILD_ID=$1
echo ${BUILD_ID}
if ! [ -x "$(command -v docker)" ]; then
  echo 'Unable to find docker command, please install Docker (https://www.docker.com/) and retry' >&2
  exit 1
fi

if [ -z "${BUILD_ID}" ]; then
  echo "Build ID is missing"
  exit 1
fi

build() {
  echo "Building image"
  docker build -t docker.io/saidsef/faas-reverse-geocoding .
  docker tag docker.io/saidsef/faas-reverse-geocoding docker.io/saidsef/faas-reverse-geocoding:${BUILD_ID}
}

push() {
  echo "Pushing image"
  docker push docker.io/saidsef/faas-reverse-geocoding
}

main() {
  build
  push
}

main
