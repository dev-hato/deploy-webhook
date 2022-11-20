#!/usr/bin/env bash

dockle_version="$(cat .dockle-version)"
curl -L -o dockle.deb "https://github.com/goodwithtech/dockle/releases/download/v${dockle_version}/dockle_${dockle_version}_Linux-64bit.deb"
sudo dpkg -i dockle.deb

docker compose -f docker-compose.yml -f "${DOCKER_COMPOSE_FILE_NAME}" pull
docker compose -f docker-compose.yml -f "${DOCKER_COMPOSE_FILE_NAME}" up -d

for image_name in $(docker compose -f docker-compose.yml -f "${DOCKER_COMPOSE_FILE_NAME}" images | awk 'OFS=":" {print $2,$3}' | tail -n +2); do
  cmd="dockle --exit-code 1 "
  cmd+="${image_name}"
  echo "> ${cmd}"
  eval "${cmd}"
done
