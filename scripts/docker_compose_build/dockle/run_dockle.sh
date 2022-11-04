#!/usr/bin/env bash

for image_name in $(docker compose -f docker-compose.yml -f "${DOCKER_COMPOSE_FILE_NAME}" images | awk 'OFS=":" {print $2,$3}' | tail -n +2); do
  cmd="dockle --exit-code 1 "
  cmd+="${image_name}"
  echo "> ${cmd}"
  eval "${cmd}"
done
