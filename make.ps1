param(
    [switch] $build,
    [switch] $release,
    [switch] $install,
    [switch] $uninstall,
    [switch] $clean,
    [switch] $installDependencies,
    [switch] $runLocal
)

$MAKE_SCRIPT_DIR = ($PSScriptRoot)
$RELEASE_DIR = "$MAKE_SCRIPT_DIR\release"

# INFO FILE
$INFO_FILE_NAME="APP_INFO.conf"
$INFO_FILE = "${RELEASE_DIR}\${INFO_FILE_NAME}"
$APP_NAME = ((Get-Content "${MAKE_SCRIPT_DIR}\${INFO_FILE_NAME}" | Select-String -Pattern '^NAME=') -split '=')[1]
$APP_VERSION = ((Get-Content "${MAKE_SCRIPT_DIR}\${INFO_FILE_NAME}" | Select-String -Pattern '^VERSION=') -split '=')[1]

# OTHERS
$BINARY = "$RELEASE_DIR\${APP_NAME}.exe"
$BINARY_LINUX="$RELEASE_DIR\${APP_NAME}"
$POWERSHELL_VENDOR_DIR = "$MAKE_SCRIPT_DIR\vendor\powershell-utils"
$PACKAGE_FILE = "$RELEASE_DIR\${APP_NAME}-${APP_VERSION}.zip"

# IMPORT LIBS
. "$POWERSHELL_VENDOR_DIR\MainUtils.ps1"

# ---------------------------------------------------------------------------- #
#                               GENERIC FUNCTIONS                              #
# ---------------------------------------------------------------------------- #
function _copyDirectory($directory, $destination, $onlyFiles) {
    if ((directoryexists "$directory")) {
        if ($onlyFiles) {
            $directory = "$directory\*"
        }
        evaladvanced "Copy-Item `"$directory`" -Destination `"$destination`" -Recurse -Force"
    } else {
        warnlog "Not found directory: $directory"
    }
}

# ---------------------------------------------------------------------------- #
#                                     MAIN                                     #
# ---------------------------------------------------------------------------- #
function _installDependencies {
    infolog "No dependencies"
}

function _install {
    $installDir = "$OTHER_APPS_DIR\$APP_NAME"
    $installBin = "$installDir\${APP_NAME}.exe"
    _installDependencies
    infolog "Install..."
    mkdir "$installDir"
    _copyDirectory "$MAKE_SCRIPT_DIR" -destination "$installDir" -onlyFiles $true
    create_shortcut_file "${APP_NAME}" "$installBin" -icon "$installDir\assets\image\logo\icon.ico" -terminal $false
    add_boot_application "$APP_NAME" "$installBin" -hidden
}

function _uninstall {
    $installDir = "$OTHER_APPS_DIR\$APP_NAME"
    infolog "Uninstall..."
    del_boot_application "$APP_NAME"
    del_shortcut_file "${APP_NAME}"
    deletedirectory "$installDir"
    deletedirectory "${home}\.config\$APP_NAME"
}

function _generatePackage() {
    Compress-Archive "$RELEASE_DIR\*" -DestinationPath "$PACKAGE_FILE" -Force
}

function _preparePackage() {
    param([bool] $isRelease)
    $releaseDate = (Get-date -Format "dd/MM/yyyy - HH:mm:ss")
    infolog "Copy necessary files..."
    _copyDirectory -directory "$MAKE_SCRIPT_DIR\assets" -destination "$RELEASE_DIR" -onlyFiles $false
    _copyDirectory -directory "$MAKE_SCRIPT_DIR\scripts" -destination "$RELEASE_DIR" -onlyFiles $false
    _copyDirectory -directory "$MAKE_SCRIPT_DIR\vendor" -destination "$RELEASE_DIR" -onlyFiles $false
    Copy-Item "$MAKE_SCRIPT_DIR\make.ps1" "$RELEASE_DIR\make.ps1"
    Copy-Item "$MAKE_SCRIPT_DIR\make.sh" "$RELEASE_DIR\make.sh"
    Copy-Item "$MAKE_SCRIPT_DIR\APP_INFO.conf" "$RELEASE_DIR\APP_INFO.conf"
    Copy-Item "$MAKE_SCRIPT_DIR\README.md" "$RELEASE_DIR\README.md"
    (Get-Content "$INFO_FILE") -replace "^RELEASE_DATE=.*", "RELEASE_DATE=$releaseDate" | Set-Content "$INFO_FILE"
    if ($isRelease) {
        _generatePackage
    }
}

function _build {
    param([bool] $isRelease)
    _clean
    mkdir "$RELEASE_DIR"
    evaladvanced "go get -u all"
    evaladvanced "go mod tidy -v"

    export GOARCH=amd64
    export GOOS=windows
    if ($isRelease) {
        infolog "Build release WINDOWS app..."
        go build -ldflags="-H windowsgui" -o "$BINARY" "$MAKE_SCRIPT_DIR\src\main.go"
    } else {
        infolog "Build WINDOWS app..."
        go build -o "$BINARY" "$MAKE_SCRIPT_DIR\src\main.go"
    }
    export GOOS=linux
    infolog "Build LINUX app..."
    go build -o "$BINARY_LINUX" "$MAKE_SCRIPT_DIR\src\main.go"
    
    _preparePackage $isRelease
}

function _clean() {
    deletedirectory "$RELEASE_DIR"
}

function main() {
    if ($clean) {
        _clean
    } elseif ($build) {
        _build
    } elseif ($release) {
        _build -isRelease $true
    } elseif ($install) {
        _uninstall
        _install
    } elseif ($uninstall) {
        _uninstall
    } elseif ($installDependencies) {
        _installDependencies
    }else {
        log "make.ps1 -build|-clean|-install|-uninstall|-installDependencies"
    }
}
main
