param(
    [switch] $build,
    [switch] $release,
    [switch] $clean,
    [switch] $run
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
$PACKAGE_FILE = "$RELEASE_DIR\${APP_NAME}-${APP_VERSION}.zip"

# ---------------------------------------------------------------------------- #
#                               GENERIC FUNCTIONS                              #
# ---------------------------------------------------------------------------- #
function _copyDirectory($directory, $destination, $onlyFiles) {
    if (!([string]::IsNullOrEmpty($directory)) -and (Test-Path -Path "$directory")) {
        if ($onlyFiles) {
            $directory = "$directory\*"
        }
        Copy-Item "$directory" -Destination "$destination" -Recurse -Force
    } else {
        Write-Host "WARN: Not found directory: $directory"
    }
}
function _deletedirectory($directory) {
    if (!([string]::IsNullOrEmpty($directory)) -and (Test-Path -Path "$directory")) {
        Remove-Item "$directory" -Recurse -Force
    }
}

# ---------------------------------------------------------------------------- #
#                                     MAIN                                     #
# ---------------------------------------------------------------------------- #
function _generatePackage() {
    Compress-Archive "$RELEASE_DIR\*" -DestinationPath "$PACKAGE_FILE" -Force
}

function _preparePackage() {
    param([bool] $isRelease)
    $releaseDate = (Get-date -Format "dd/MM/yyyy - HH:mm:ss")
    Write-Host "INFO: Copy necessary files..."
    _copyDirectory -directory "$MAKE_SCRIPT_DIR\assets" -destination "$RELEASE_DIR" -onlyFiles $false
    _copyDirectory -directory "$MAKE_SCRIPT_DIR\scripts" -destination "$RELEASE_DIR" -onlyFiles $false
    _deletedirectory "$RELEASE_DIR/vendor/golangutils"
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
    New-Item -Path "$RELEASE_DIR" -ItemType Directory -Force | Out-Null
    go get -u all
    go mod tidy -v

    set-item -force -path "env:GOARCH" -value "amd64"
    set-item -force -path "env:GOOS" -value "windows"
    if ($isRelease) {
        Write-Host "INFO: Build release WINDOWS app..."
        go build -ldflags="-H windowsgui" -o "$BINARY" "$MAKE_SCRIPT_DIR\src\main.go"
    } else {
        Write-Host "INFO: Build WINDOWS app..."
        go build -o "$BINARY" "$MAKE_SCRIPT_DIR\src\main.go"
    }
    set-item -force -path "env:GOOS" -value "linux"
    Write-Host "INFO: Build LINUX app..."
    go build -o "$BINARY_LINUX" "$MAKE_SCRIPT_DIR\src\main.go"

    _preparePackage $isRelease
}

function _clean() {
    _deletedirectory "$RELEASE_DIR"
}

function main() {
    if ($clean) {
        _clean
    } elseif ($build) {
        _build
    } elseif ($release) {
        _build -isRelease $true
    } elseif ($run) {
        Push-Location
        Set-Location "$RELEASE_DIR"
        Start-Process "$BINARY" -Wait
        Pop-Location
    } elseif ($installDependencies) {
        _installDependencies
    }else {
        log "make.ps1 -build|-clean|-release|-run"
    }
}
main
