#!/usr/bin/env bash

TAG_NAME="${1}"
export TAG_NAME

cd "${LOCAL_REPO_PATH}" || exit 1
docker compose pull
docker compose up -d --wait
