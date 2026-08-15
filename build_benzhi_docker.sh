#!/usr/bin/env sh
set -eu

platform="${1:-linux/amd64}"
image="expense-review-cli:benzhi"

docker build --platform "$platform" -f benzhi.Dockerfile -t "$image" .
docker run --rm --platform "$platform" "$image" go build ./...
docker run --rm --platform "$platform" "$image"
