#!/usr/bin/env bash

TAG_NAME="${1}"
export TAG_NAME

cd /repo || exit 1
docker compose -f docker-compose.yml -f onpremises.docker-compose.yml pull
docker compose -f docker-compose.yml -f onpremises.docker-compose.yml up -d --wait
