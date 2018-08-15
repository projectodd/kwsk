#!/bin/bash

set -e

IMAGE_PREFIX=${IMAGE_PREFIX:-"projectodd"}
IMAGE_TAG=${IMAGE_TAG:-"latest"}
IMAGE_PUSH=${IMAGE_PUSH:-"false"}

pushd cmd/kwsk-runtime-shim
go build
CGO_ENABLED=0 go build -a -o kwsk-runtime-shim.cgo
docker build -f Dockerfile.nodejs6 . --tag ${IMAGE_PREFIX}/kwsk-nodejs6action:${IMAGE_TAG}
docker build -f Dockerfile.nodejs8 . --tag ${IMAGE_PREFIX}/kwsk-action-nodejs-v8:${IMAGE_TAG}
docker build -f Dockerfile.python2 . --tag ${IMAGE_PREFIX}/kwsk-python2action:${IMAGE_TAG}
docker build -f Dockerfile.python3 . --tag ${IMAGE_PREFIX}/kwsk-python3action:${IMAGE_TAG}
docker build -f Dockerfile.java8   . --tag ${IMAGE_PREFIX}/kwsk-java8action:${IMAGE_TAG}
docker images
if [ "$IMAGE_PUSH" == "true" ]; then
  echo "$DOCKER_PASSWORD" | docker login -u "${DOCKER_USER}" --password-stdin
  docker push ${IMAGE_PREFIX}/kwsk-nodejs6action:${IMAGE_TAG}
  docker push ${IMAGE_PREFIX}/kwsk-action-nodejs-v8:${IMAGE_TAG}
  docker push ${IMAGE_PREFIX}/kwsk-python2action:${IMAGE_TAG}
  docker push ${IMAGE_PREFIX}/kwsk-python2action:${IMAGE_TAG}
  docker push ${IMAGE_PREFIX}/kwsk-java8action:${IMAGE_TAG}
fi
popd
