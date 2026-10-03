#!/bin/bash

APP_NAME="$1"
NO_CONSOLE_WINDOW="$2"
ROOT_DIR="$PWD"
BINARY_DIR="$ROOT_DIR/bin"
GO_RELEASE_FLAGS="-s -w"
GO_RELEASE_NO_CONSOLE_FLAGS=""

if [[ "$NO_CONSOLE_WINDOW" == "1" ]]; then
	GO_RELEASE_NO_CONSOLE_FLAGS="-H=windowsgui"
fi

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
go build -ldflags="$GO_RELEASE_FLAGS $GO_RELEASE_NO_CONSOLE_FLAGS" -o "$BINARY_DIR/${APP_NAME}.exe" "$ROOT_DIR"
