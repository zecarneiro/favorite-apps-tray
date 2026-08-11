#!/bin/bash

APP_NAME="$1"
ROOT_DIR="$PWD"
BINARY_DIR="$ROOT_DIR/bin"
GO_RELEASE_FLAGS="-s -w"

echo ">>> Create Binary directory: $BINARY_DIR"
rm -rf "$BINARY_DIR"
mkdir -p "$BINARY_DIR"

echo ">>> Building for linux..."
export GOOS=linux
export GOARCH=amd64
go build -ldflags="$GO_RELEASE_FLAGS" -o "$BINARY_DIR/$APP_NAME" "$ROOT_DIR"
echo ">>> Building for windows..."
export GOOS=windows
export GOARCH=amd64
go build -ldflags="$GO_RELEASE_FLAGS"  -o "$BINARY_DIR/${APP_NAME}.exe" "$ROOT_DIR"
