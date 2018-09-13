#!/bin/bash

TAG=$TRAVIS_BUILD_NUMBER-$TRAVIS_COMMIT

docker build -t pitaya-admin .

docker login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"

docker tag pitaya-admin:latest tfgco/pitaya-admin:$TAG
docker push tfgco/pitaya-admin:$TAG
docker tag pitaya-admin:latest tfgco/pitaya-admin:latest
docker push tfgco/pitaya-admin

DOCKERHUB_LATEST=$(python ./scripts/get_latest_tag.py)

if [ "$DOCKERHUB_LATEST" != "$TAG" ]; then
    echo "Last version is not in docker hub!"
    echo "docker hub: $DOCKERHUB_LATEST, expected: $TAG"
    exit 1
fi
