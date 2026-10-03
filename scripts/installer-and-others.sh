#!/bin/bash

ROOT_DIR="$1"
RELEASE_DIR="$2"
DEPLOY_DIR="$3"
INSTALLERS_DIR="$ROOT_DIR/scripts/installers"
OTHERS_SCRIPTS_DIR="$ROOT_DIR/scripts/others"
RELEASE_SCRIPTS_DIR="$DEPLOY_DIR/scripts"
IMAGES_DIR="$ROOT_DIR/docs/images"

echo ">>> Copy images..."
cp "$IMAGES_DIR/logo/linux.png" "$DEPLOY_DIR/linux.png"
cp "$IMAGES_DIR/logo/win.ico" "$DEPLOY_DIR/win.ico"

echo ">>> Copy installer and uninstaller..."
SCOOP_INSTALLER="$RELEASE_DIR/$APP_NAME.json"
LINUX_INSTALLER="$RELEASE_DIR/${APP_NAME}-install.sh"
LINUX_UNINSTALLER="$RELEASE_DIR/${APP_NAME}-uninstall.sh"
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
