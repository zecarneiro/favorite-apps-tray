#!/bin/bash
# shellcheck disable=SC2155

declare appName="$1"
declare typeInfo="$2"


declare SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
declare TRAY_APP_CONFIG_DIR="${HOME}/.config/${appName}"
declare JSON_FILE="${TRAY_APP_CONFIG_DIR}/apps-info-${typeInfo}.json"
declare APPS_ARR="[]"

function saveToFile {
    if [ -f "$JSON_FILE" ]; then
        rm "$JSON_FILE"
    fi 
    echo "INFO: ${appName} - Generating: $JSON_FILE"
    echo "$APPS_ARR" | tee "$JSON_FILE" >/dev/null
}

function main {
    if [[ "${typeInfo}" == "shortcuts" ]]; then
        APPS_ARR=$(gjs "${SCRIPT_DIR}/desktop-file-info-in-gnome.js")
        saveToFile
    fi
}
main
