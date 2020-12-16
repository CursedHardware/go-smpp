#!/bin/bash
set -xeuo pipefail

mkdir -p dist

function build() {
  local IFS='/'
  local NAME="smpp-$1"
  for i in "${@:2}"; do
    read -r -a platform <<<"$i"
    local GOOS="${platform[0]}"
    local GOARCH="${platform[1]}"
    env GOOS="$GOOS" GOARCH="$GOARCH" \
      go build \
      -ldflags '-s -w' \
      -o "dist/$NAME-$GOOS-$GOARCH" \
      "./cmd/$NAME"
  done
}

build receiver linux/amd64 linux/arm linux/arm64
build repl linux/amd64 linux/arm linux/arm64 darwin/amd64

upx dist/*
