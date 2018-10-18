#!/bin/sh

export GOOS=linux
export GOARCH=arm
export CGO_ENABLED=0

set -e
set -x

go build -o release/linux/arm/drone-docker \
    github.com/drone/drone-docker
