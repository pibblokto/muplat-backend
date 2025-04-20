#!/bin/bash


TAG=shittiest
if [[ $1 ]]; then
    TAG=$1
fi

docker build -t europe-west3-docker.pkg.dev/piblokto/muplat/server:${TAG} .
docker push europe-west3-docker.pkg.dev/piblokto/muplat/server:${TAG}
