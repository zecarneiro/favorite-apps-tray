#!/bin/bash
# AUTHOR: Jose M. Noronha

declare OTHERAPPS_DIR="$HOME/.otherapps"
declare INSTALL_DIR="$OTHERAPPS_DIR/{APP_NAME}"
declare SHORTCUT="$HOME/.local/share/applications/{APP_NAME}.desktop"
declare SYMBOLIC_SYSTEM_FILE="/usr/bin/{APP_NAME}"

function _printInfo() {
    local operation="$1"
    local message="$2"
    echo "[INFO] ${operation}: ${message}"
}

function _printError() {
    local message="$1"
    echo "[ERROR] ${message}"
}

function _uninstall() {
    if [ -h "$SYMBOLIC_SYSTEM_FILE" ]; then
        _printInfo "Remove" "$SYMBOLIC_SYSTEM_FILE"
        sudo rm "$SYMBOLIC_SYSTEM_FILE"
    fi
    if [ -d "$INSTALL_DIR" ]; then
        _printInfo "Remove" "$INSTALL_DIR"
        rm -rf "$INSTALL_DIR"
    fi
    if [ -f "$SHORTCUT" ]; then
        _printInfo "Remove" "$SHORTCUT"
        rm "$SHORTCUT"
    fi
}
_uninstall
