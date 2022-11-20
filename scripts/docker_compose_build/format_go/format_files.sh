#!/usr/bin/env bash

echo "$(go env GOPATH)/bin" >>"${GITHUB_PATH}"
go install golang.org/x/tools/cmd/goimports

go_version="$(go version | awk '{print $3}' | awk -F . 'BEGIN {OFS="."} {print $1,$2}' | sed -e 's/go//g')"
go mod tidy -go="${go_version}"
goimports -l -w .
