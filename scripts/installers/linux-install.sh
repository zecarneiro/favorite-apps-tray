#!/bin/bash

declare OTHERAPPS_DIR="$HOME/.local/opt"
declare INSTALL_DIR="$OTHERAPPS_DIR/{APP_NAME}"
declare SHORTCUT="$HOME/.local/share/applications/{APP_NAME}.desktop"
declare ZIP_FILE="$OTHERAPPS_DIR/favorite-apps-tray-{APP_VERSION}.zip"
declare SYMBOLIC_SYSTEM_FILE="/usr/bin/{APP_NAME}"
declare URL="https://github.com/zecarneiro/favorite-apps-tray/releases/download/v{APP_VERSION}/favorite-apps-tray-{APP_VERSION}.zip"

function _printInfo() {
    local operation="$1"
    local message="$2"
    echo "[INFO] ${operation}: ${message}"
}

function _printError() {
    local message="$1"
    echo "[ERROR] ${message}"
}

function _check_dependencies() {
    if [ ! "$(command -v wget)" ]; then
        _printError "Please install wget!"
        exit 1
    fi
    if [ ! "$(command -v unzip)" ]; then
        _printError "Please install unzip!"
        exit 1
    fi
    if [ ! "$(command -v zenity)" ]; then
        _printError "Please install zenity!"
        exit 1
    fi
    if [ ! "$(command -v notify-send)" ]&&[ ! "$(command -v kdialog)" ]; then
        _printError "Please install notify-send or kdialog!"
        exit 1
    fi
    if [ ! "$(command -v inkscape)" ]; then
        _printError "Please install inkscape!"
        exit 1
    fi
}

function _install() {
    local data="[Desktop Entry]
Version=1.0
Type=Application
Terminal=true
Exec=$INSTALL_DIR/{APP_NAME}
Name={APP_DISPLAY_NAME}
Comment={APP_DISPLAY_NAME}
Icon=$INSTALL_DIR/linux.png"

    # Start
    _check_dependencies
    if [ -f "$ZIP_FILE" ]; then
        _printInfo "Remove" "$ZIP_FILE"
        rm "$ZIP_FILE"
    fi

    _printInfo "Create" "$INSTALL_DIR"
    mkdir -p "$INSTALL_DIR"

    _printInfo "Download" "{APP_DISPLAY_NAME}"
    wget -O "$ZIP_FILE" "$URL" -q --show-progress || exit 1

    _printInfo "Install" "{APP_DISPLAY_NAME}"
    unzip -o "$ZIP_FILE" -d "$INSTALL_DIR" || exit 1
    echo -e "$data" | tee "$SHORTCUT" >/dev/null
    chmod +x "$SHORTCUT" || exit 1
    chmod +x "$INSTALL_DIR/{APP_NAME}" || exit 1
    if [ -f "$ZIP_FILE" ]; then
        _printInfo "Remove" "$ZIP_FILE"
        rm "$ZIP_FILE"
    fi
    if [ -f "$SYMBOLIC_SYSTEM_FILE" ]; then
        sudo rm "$SYMBOLIC_SYSTEM_FILE"
    fi
    sudo ln "$INSTALL_DIR/{APP_NAME}" "$SYMBOLIC_SYSTEM_FILE"
}
_install
