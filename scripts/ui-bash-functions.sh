#!/bin/bash
# Author: José M. C. Noronha

alias zenity="zenity 2>/dev/null"

function notify {
$appId = "${env:FAT_NOTIFY_APP_ID}"
$title = "${env:FAT_NOTIFY_TITLE}"
$message = "${env:FAT_NOTIFY_MESSAGE}"
$icon = "${env:FAT_NOTIFY_ICON}"


    local appId="$1"
    local title="$2"
    local message="$3"
    local icon="$4"
    if [[ -n "${icon}" ]]; then
        notify-send -a "$appId" -i "$icon" -t 60000 "$title" "$message"
    else
        notify-send -a "$appId" -t 60000 "$title" "$message"
    fi
}

function oknotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Success" "$message" "$icon"
}

function infonotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Information" "$message" "$icon"
}

function warnnotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Warning" "$message" "$icon"
}

function errornotify {
    local appId="$1"
    local message="$2"
    local icon="$3"
    notify "$appId" "Error" "$message" "$icon"
}

function selectfiledialog {
    local selectedFile=$(zenity --file-selection --modal)
    echo "{ \"selected\": \"$selectedFile\" }"
}

function svg_to_png {
    local svg_file="$1"
    local png_file="$2"
    if [[ -f "$svg_file" ]]; then
        local png_dir="$(dirname "$png_file")"
        if [[ ! -d "$png_dir" ]]; then
            mkdir "$png_dir"
        fi
        eval "inkscape \"$svg_file\" -o \"$png_file\" --export-overwrite -w 32 -h 32"
    fi
}
