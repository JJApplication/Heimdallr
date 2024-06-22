#!/usr/bin/env bash

# build for app
echo "GOPROXY: $(go env GOPROXY)"
echo "GOOS: $(go env GOOS)"
echo "GOARCH: $(go env GOARCH)"

date -u

echo "start to build"
go build -mod=mod -ldflags='-w -s' -trimpath -o heimdallr ./
echo "done"