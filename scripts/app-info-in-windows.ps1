param([string] $appName, [string] $typeInfo)

$SCRIPT_DIR = ($PSScriptRoot)
$TRAY_APP_CONFIG_DIR = "${home}\.config\${appName}"
$JSON_FILE = "${TRAY_APP_CONFIG_DIR}\apps-info-${typeInfo}.json"
$APPS_ARR = @()

# IMPORT LIBS
. "$SCRIPT_DIR\..\vendor\powershell-utils\MainUtils.ps1"

function saveToFile {
    deletefile "$JSON_FILE"
    infolog "${appName} - Generating: $JSON_FILE"
    (ConvertTo-Json $APPS_ARR -Depth 2 | Out-String) | teecustom "$JSON_FILE"
}

function main {
    mkdir "$TRAY_APP_CONFIG_DIR"
    if ("windows-apps" -eq $typeInfo) {
        (New-Object -ComObject Shell.Application).NameSpace('shell:AppsFolder').Items() | ForEach-Object {
            $nameApp = $_.Name
            $pathApp = $_.Path
            $path = "shell:AppsFolder\'$pathApp'"
            $appInfo = @{
                displayName="$nameApp";
                command="Start-Process -FilePath $path";
                shortcut="$nameApp"
            }
            $APPS_ARR += $appInfo
        }
        saveToFile
    } elseif ("shortcuts" -eq $typeInfo) {
        $menuDirs = "C:\ProgramData\Microsoft\Windows\Start Menu\Programs", "$home\AppData\Roaming\Microsoft\Windows\Start Menu"
        $menuDirs | ForEach-Object {
            $menuDir = $_
            Get-ChildItem -Path "$menuDir" -Include "*.lnk" -File -Recurse | ForEach-Object {
                $shortcut = $_
                $appInfo = [PSCustomObject]@{
                    displayName=filename(basename("$shortcut"));
                    command="Start-Process -FilePath '$shortcut'";
                    shortcut=basename("$shortcut")
                    icon="$shortcut"
                }
                $APPS_ARR += $appInfo
            }
        }        
        saveToFile
    }
}
main
