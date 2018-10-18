#!/bin/sh

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

set -e
set -x

go build -o release/linux/amd64/drone-docker \
    github.com/drone/drone-docker
