#!/bin/bash
# AUTHOR: José M. Noronha

declare MAKE_SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
declare RELEASE_DIR="$MAKE_SCRIPT_DIR/release"

# INFO FILE
declare INFO_FILE_NAME="APP_INFO.conf"
declare INFO_FILE="${RELEASE_DIR}/${INFO_FILE_NAME}"
declare APP_NAME="$(cat "${MAKE_SCRIPT_DIR}/${INFO_FILE_NAME}" | grep "NAME" | cut -d'=' -f2)"
declare APP_VERSION="$(cat "${MAKE_SCRIPT_DIR}/${INFO_FILE_NAME}" | grep "VERSION" | cut -d'=' -f2)"

# OTHERS
declare BINARY="$RELEASE_DIR/${APP_NAME}"
declare BINARY_WINDOWS="$RELEASE_DIR/${APP_NAME}.exe"
declare BASH_VENDOR_DIR="$MAKE_SCRIPT_DIR/vendor/bash-utils"
declare PACKAGE_FILE_NAME="${APP_NAME}-${APP_VERSION}.zip"
declare PACKAGE_FILE="$MAKE_SCRIPT_DIR/${PACKAGE_FILE_NAME}"

# IMPORT LIBS
. "$BASH_VENDOR_DIR/main-utils.sh"

# ---------------------------------------------------------------------------- #
#                               GENERIC FUNCTIONS                              #
# ---------------------------------------------------------------------------- #
function _copyDirectory {
    local directory="$1"
    local destination="$2"
    local onlyFiles="$3"
    if [ $(directoryexists "$directory") == true ]; then
        if [ $onlyFiles == true ]; then
            directory="$directory/."
        fi
        evaladvanced "cp -r \"$directory\" \"$destination\""
    else
        warnlog "Not found directory: $directory"
    fi
}

# ---------------------------------------------------------------------------- #
#                                     MAIN                                     #
# ---------------------------------------------------------------------------- #
function _installDependencies() {
    local package_list=(zenity libnotify-bin)
    infolog "Install all dependencies"
    for package_name in "${package_list[@]}"; do
        evaladvanced "sudo apt install $package_name -y"
    done
}

function _install() {
    local installDir="$OTHER_APPS_DIR/$APP_NAME"
    local installBin="$installDir/${APP_NAME}"
    infolog "Install..."
    _installDependencies
    mkdir "$installDir"
    _copyDirectory "$MAKE_SCRIPT_DIR" "$installDir" true
    create_shortcut_file --name "${APP_NAME}" --exec "$installBin" --icon "$installDir/assets/image/logo/icon.png"
    add_boot_application "$APP_NAME" "$installBin" 1
    chmod -R 777 "$installDir"
    evaladvanced "sudo update-desktop-database"
}

function _uninstall() {
    local installDir="$OTHER_APPS_DIR/$APP_NAME"
    infolog "Uninstall..."
    del_boot_application "$APP_NAME"
    del_shortcut_file "${APP_NAME}"
    deletedirectory "$installDir"
    deletedirectory "${HOME}/.config/$APP_NAME"
    evaladvanced "sudo update-desktop-database"
}

function _generatePackage() {
    pushd .
    cd "$RELEASE_DIR"
    zip -rq "$PACKAGE_FILE" .
    popd
    mv "$PACKAGE_FILE" "$RELEASE_DIR/$PACKAGE_FILE_NAME"
}

function _preparePackage() {
    local isRelease="$1"
    local releaseDate=$(date '+%d/%m/%Y %H:%M:%S')
    infolog "Copy necessary files..."
    _copyDirectory "$MAKE_SCRIPT_DIR/assets" "$RELEASE_DIR" false
    _copyDirectory "$MAKE_SCRIPT_DIR/scripts" "$RELEASE_DIR" false
    _copyDirectory "$MAKE_SCRIPT_DIR/vendor" "$RELEASE_DIR" false
    deletedirectory "$RELEASE_DIR/vendor/golangutils"
    cp "$MAKE_SCRIPT_DIR/make.ps1" "$RELEASE_DIR/make.ps1"
    cp "$MAKE_SCRIPT_DIR/make.sh" "$RELEASE_DIR/make.sh"
    cp "$MAKE_SCRIPT_DIR/APP_INFO.conf" "$RELEASE_DIR/APP_INFO.conf"
    cp "$MAKE_SCRIPT_DIR/README.md" "$RELEASE_DIR/README.md"
    sed -i "s#RELEASE_DATE=.*#RELEASE_DATE=${releaseDate}#g" "$INFO_FILE"
    if [[ "${isRelease}" == "true" ]]; then
        _generatePackage
    fi
}

function _build() {
    local isRelease="$1"
    _clean
    mkdir "$RELEASE_DIR"
    evaladvanced "go get -u all"
    evaladvanced "go mod tidy -v"

    export GOARCH=amd64
    export GOOS=windows
    if [[ "${isRelease}" == "true" ]]; then
        infolog "Build release WINDOWS app..."
        go build -ldflags="-H windowsgui" -o "$BINARY_WINDOWS" "$MAKE_SCRIPT_DIR/src/main.go"
    else
        infolog "Build WINDOWS app..."
        go build -o "$BINARY_WINDOWS" "$MAKE_SCRIPT_DIR/src/main.go"
    fi
    export GOOS=linux
    infolog "Build release LINUX app..."
    go build -o "$BINARY" "$MAKE_SCRIPT_DIR/src/main.go"

    _preparePackage $isRelease
}

function _clean() {
    deletedirectory "$RELEASE_DIR"
    deletefile "$PACKAGE_FILE"
}

function main() {
    case "${1}" in
        -clean) _clean ;;
        -build) _build ;;
        -release) _build true ;;
        -install)
            _uninstall
            _install
        ;;
        -uninstall) _uninstall ;;
        -installDependencies) _installDependencies ;;
        *) log "make.ps1 -build|-clean|-install|-uninstall|-installDependencies" ;;
    esac
}
main "$@"
