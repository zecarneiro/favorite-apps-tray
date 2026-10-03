#!/bin/bash

SO_TYPE="$1"; shift
APP_NAME="$1"; shift
APP_VERSION="$1"; shift
APP_DISPLAY_NAME="$*"
ROOT_DIR="$PWD"
RELEASE_DIR="$ROOT_DIR/release"
DEPLOY_DIR="$RELEASE_DIR/$APP_NAME"
BINARY_DIR="$ROOT_DIR/bin"
CONF_FILE="$DEPLOY_DIR/APP_INFO.conf"
SCRIPTS_DIR="$ROOT_DIR/scripts"

echo ">>> Create release directory..."
rm -rf "$RELEASE_DIR"
mkdir -p "$DEPLOY_DIR"

echo ">>> Copy binaries..."
cp "$BINARY_DIR/$APP_NAME" "$DEPLOY_DIR/$APP_NAME"
cp "$BINARY_DIR/${APP_NAME}.exe" "$DEPLOY_DIR/${APP_NAME}.exe"

echo ">>> Generate Conf file"
echo "NAME=${APP_NAME}" | tee -a "${CONF_FILE}"
echo "DISPLAY_NAME=${APP_DISPLAY_NAME}" | tee -a "${CONF_FILE}"
echo "VERSION=${APP_VERSION}" | tee -a "${CONF_FILE}"
echo "RELEASE=true" | tee -a "${CONF_FILE}"
echo "RELEASE_DATE=$(date '+%d/%m/%Y %H:%M:%S')" | tee -a "${CONF_FILE}"

# Process Others
. "$SCRIPTS_DIR/installer-and-others.sh" "$ROOT_DIR" "$RELEASE_DIR" "$DEPLOY_DIR"


echo ">>> Generate package file..."
PACKAGE_FILE=""
if [[ "$SO_TYPE" == "windows" ]]; then
    PACKAGE_FILE="$RELEASE_DIR\${APP_NAME}-${APP_VERSION}.zip"
    powershell.exe -Command "Compress-Archive '$DEPLOY_DIR\*' -DestinationPath '$PACKAGE_FILE'" -Force
else
    PACKAGE_FILE=$RELEASE_DIR/${APP_NAME}-${APP_VERSION}.zip
    cd "$DEPLOY_DIR" || exit 1
    zip -rq "$PACKAGE_FILE" .
    cd "$ROOT_DIR" || exit 1
fi

echo ">>> Delete unnecessary files and directories..."
rm -rf "$BINARY_DIR"
