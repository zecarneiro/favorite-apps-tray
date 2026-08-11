#!/bin/bash

APP_NAME="$1"
ROOT_DIR="$PWD"
RELEASE_DIR="$ROOT_DIR/release"
DEPLOY_DIR="$ROOT_DIR/$APP_NAME"
BINARY_DIR="$ROOT_DIR/bin"

echo ">>> Cleanning..."
_delete_dir() {
    echo ">>> Delete directory: $1"
    rm -rf "$1"
}
_delete_dir "$RELEASE_DIR"
_delete_dir "$DEPLOY_DIR"
_delete_dir "$BINARY_DIR"
