#!/usr/bin/env bash

set -eu
set -o pipefail

go get

version=$(cat VERSION || echo "no_version_info")
commit=$(git rev-parse --short HEAD || echo "not_git_repository")
buildTimestamp=$(date --iso-8601=seconds)

gofmt -w *.go

go build -ldflags "-X 'main.version=$version' -X 'main.commit=$commit' -X 'main.buildTimestamp=$buildTimestamp'"
