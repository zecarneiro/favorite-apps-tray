#!/bin/bash

SO_TYPE="$1"; shift
APP_NAME="$1"; shift
APP_VERSION="$1"; shift
APP_DISPLAY_NAME="$*"
ROOT_DIR="$PWD"
RELEASE_DIR="$ROOT_DIR/release"
DEPLOY_DIR="$ROOT_DIR/$APP_NAME"
BINARY_DIR="$ROOT_DIR/bin"
IMAGES_DIR="$ROOT_DIR/docs/images"
INSTALLERS_DIR="$ROOT_DIR/scripts/installers"
OTHERS_SCRIPTS_DIR="$ROOT_DIR/scripts/others"
RELEASE_SCRIPTS_DIR="$RELEASE_DIR/scripts"
CONF_FILE="$RELEASE_DIR/APP_INFO.conf"
echo ">>> Create release directory: $RELEASE_DIR"
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

echo ">>> Copy binaries..."
cp "$BINARY_DIR/$APP_NAME" "$RELEASE_DIR/$APP_NAME"
cp "$BINARY_DIR/${APP_NAME}.exe" "$RELEASE_DIR/${APP_NAME}.exe"

echo ">>> Copy images..."
cp "$IMAGES_DIR/logo/linux.png" "$RELEASE_DIR/linux.png"
cp "$IMAGES_DIR/logo/win.ico" "$RELEASE_DIR/win.ico"

echo ">>> Create deploy directory: $DEPLOY_DIR"
rm -rf "$DEPLOY_DIR"
mkdir -p "$DEPLOY_DIR"

echo ">>> Generate Conf file"
echo "NAME=${APP_NAME}" | tee -a "${CONF_FILE}"
echo "DISPLAY_NAME=${APP_DISPLAY_NAME}" | tee -a "${CONF_FILE}"
echo "VERSION=${APP_VERSION}" | tee -a "${CONF_FILE}"
echo "RELEASE=true" | tee -a "${CONF_FILE}"
echo "RELEASE_DATE=$(date '+%d/%m/%Y %H:%M:%S')" | tee -a "${CONF_FILE}"

echo ">>> Copy installer and uninstaller..."
SCOOP_INSTALLER="$DEPLOY_DIR/$APP_NAME.json"
LINUX_INSTALLER="$DEPLOY_DIR/${APP_NAME}-install.sh"
LINUX_UNINSTALLER="$DEPLOY_DIR/${APP_NAME}-uninstall.sh"
cp "$INSTALLERS_DIR/scoop.json" "$SCOOP_INSTALLER"
cp "$INSTALLERS_DIR/linux-install.sh" "$LINUX_INSTALLER"
cp "$INSTALLERS_DIR/linux-uninstall.sh" "$LINUX_UNINSTALLER"
declare -A REPLACER=(
    [{APP_VERSION}]="${APP_VERSION}"
    [{APP_NAME}]="${APP_NAME}"
    [{APP_DISPLAY_NAME}]="${APP_DISPLAY_NAME}"
)
for key in "${!REPLACER[@]}"; do
    sed -i "s#$key#${REPLACER[$key]}#g" "$SCOOP_INSTALLER"
    sed -i "s#$key#${REPLACER[$key]}#g" "$LINUX_INSTALLER"
    sed -i "s#$key#${REPLACER[$key]}#g" "$LINUX_UNINSTALLER"
done

echo ">>> Copy other files..."
mkdir -p "$RELEASE_SCRIPTS_DIR"
cp "$OTHERS_SCRIPTS_DIR/app-info-for-gnome.js" "$RELEASE_SCRIPTS_DIR/app-info-for-gnome.js"
cp "$OTHERS_SCRIPTS_DIR/app-info-for-windows.ps1" "$RELEASE_SCRIPTS_DIR/app-info-for-windows.ps1"
cp "$OTHERS_SCRIPTS_DIR/icon-extractor-for-windows.ps1" "$RELEASE_SCRIPTS_DIR/icon-extractor-for-windows.ps1"

echo ">>> Generate package file..."
if [[ "$SO_TYPE" == "windows" ]]; then
    powershell.exe -Command "Compress-Archive '$RELEASE_DIR\*' -DestinationPath '$DEPLOY_DIR\${APP_NAME}-${APP_VERSION}.zip'" -Force
else
    cd "$RELEASE_DIR" || exit 1
    zip -rq "$DEPLOY_DIR/${APP_NAME}-${APP_VERSION}.zip" .
    cd "$ROOT_DIR" || exit 1
fi

echo ">>> Delete unnecessary files and directories..."
rm -rf "$RELEASE_DIR"
