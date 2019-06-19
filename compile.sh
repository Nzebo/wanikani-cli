#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        echo "Building $package-$GOOS-$GOARCH"
        go build -v -o $package-$GOOS-$GOARCH $package
    done
done