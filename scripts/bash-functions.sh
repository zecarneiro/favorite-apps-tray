#!/bin/bash
# Author: José M. C. Noronha
# shellcheck disable=SC2155

# ---------------------------------------------------------------------------- #
#                                 APP INFO AREA                                #
# ---------------------------------------------------------------------------- #
declare appName="$FAT_INFO_APP_NAME"
declare typeInfo="$FAT_INFO_TYPE"
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

function appinfo {
    if [[ "${typeInfo}" == "shortcuts" ]]; then
        APPS_ARR=$(gjs "${SCRIPT_DIR}/desktop-file-info-in-gnome.js")
        saveToFile
    fi
}

# ---------------------------------------------------------------------------- #
#                                    UI AREA                                   #
# ---------------------------------------------------------------------------- #
alias zenity="zenity 2>/dev/null"

function notify {
    local appId="$FAT_NOTIFY_APP_ID"
    local title="$FAT_NOTIFY_TITLE"
    local message="$FAT_NOTIFY_MESSAGE"
    local icon="$FAT_NOTIFY_ICON"
    if [[ -n "${icon}" ]]; then
        notify-send -a "$appId" -i "$icon" -t 60000 "$title" "$message"
    else
        notify-send -a "$appId" -t 60000 "$title" "$message"
    fi
}

function selectfiledialog {
    local selectedFile=$(zenity --file-selection --modal)
    echo "{ \"selected\": \"$selectedFile\" }"
}

function messagedialog {
    local title="$FAT_DIALOG_TITLE"
    local message="$FAT_DIALOG_MESSAGE"
    local icon="$FAT_DIALOG_ICON"
    zenity --info --text="$message" --title="$title" --icon="$icon"
}

# ---------------------------------------------------------------------------- #
#                                  OTHERS AREA                                 #
# ---------------------------------------------------------------------------- #
function iconextractor {
    local svg_file="$FAT_ICON_EXTRACTOR_FILE"
    local png_file="$FAT_ICON_EXTRACTOR_DEST"
    if [[ -f "$svg_file" ]]; then
        local png_dir="$(dirname "$png_file")"
        if [[ ! -d "$png_dir" ]]; then
            mkdir "$png_dir"
        fi
        eval "inkscape \"$svg_file\" -o \"$png_file\" --export-overwrite -w 32 -h 32"
    fi
}
