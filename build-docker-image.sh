#!/bin/sh
DIRBASE=$(dirname $0)

# Build static binary for Linux image
CGO_ENABLED=0 GOOS=linux go build -a -o "${DIRBASE}"/build/urlinfo github.com/mayo/urlinfo/cmd/urlinfo

# Build a minimal Docker image for the service
docker build -t urlinfo -f "${DIRBASE}"/Dockerfile "${DIRBASE}"/