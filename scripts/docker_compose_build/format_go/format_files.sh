#!/usr/bin/env bash

go_version="$(go version | awk '{print $3}' | awk -F . 'BEGIN {OFS="."} {print $1,$2}' | sed -e 's/go//g')"
go mod tidy -go="${go_version}"
goimports -l -w .
