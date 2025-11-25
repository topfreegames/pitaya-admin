#!/bin/bash

COMMIT=$TRAVIS_COMMIT
TAG=$TRAVIS_BUILD_NUMBER-${COMMIT:0:7}

# Setup buildx for multi-platform builds
docker buildx create --use --name multiarch-builder --driver docker-container || docker buildx use multiarch-builder

docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

# Build and push multi-platform images (amd64 and arm64)
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --tag tfgco/pitaya-admin:$TAG \
  --tag tfgco/pitaya-admin:latest \
  --push \
  .

DOCKERHUB_LATEST=$(python ./scripts/get_latest_tag.py)

if [ "$DOCKERHUB_LATEST" != "$TAG" ]; then
    echo "Last version is not in docker hub!"
    echo "docker hub: $DOCKERHUB_LATEST, expected: $TAG"
    exit 1
fi
