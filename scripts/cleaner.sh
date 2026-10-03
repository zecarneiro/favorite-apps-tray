#!/bin/bash

ROOT_DIR="$PWD"
RELEASE_DIR="$ROOT_DIR/release"
BINARY_DIR="$ROOT_DIR/bin"

echo ">>> Cleanning..."
_delete_dir() {
    echo ">>> Delete directory: $1"
    rm -rf "$1"
}
_delete_dir "$RELEASE_DIR"
_delete_dir "$BINARY_DIR"
