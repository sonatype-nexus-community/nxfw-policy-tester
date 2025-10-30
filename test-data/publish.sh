#!/bin/bash

set -e

for dockerfile in Dockerfile.*; do
    suffix="${dockerfile#Dockerfile.}"
    tag=$(echo "$suffix" | sed 's/\([A-Z]\)/-\1/g' | sed 's/^-//')
    docker buildx build --platform linux/amd64,linux/arm64 -f "$dockerfile" -t "sonatypecommunity/docker-policy-demo:$tag" --push .
done

