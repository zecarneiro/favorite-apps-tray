$typeInfo = "${env:FAT_INFO_TYPE}"
$APPS_ARR = @()

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
} elseif ("shortcuts" -eq $typeInfo) {
    $menuDirs = "C:\ProgramData\Microsoft\Windows\Start Menu\Programs", "$home\AppData\Roaming\Microsoft\Windows\Start Menu"
    $menuDirs | ForEach-Object {
        $menuDir = $_
        Get-ChildItem -Path "$menuDir" -Include "*.lnk" -File -Recurse | ForEach-Object {
            $shortcut = $_
            $basename = ([System.IO.Path]::GetFileName("$shortcut"))
            $filename = ([System.IO.Path]::GetFileNameWithoutExtension("$basename"))
            $appInfo = [PSCustomObject]@{
                displayName="$filename";
                command="Start-Process -FilePath '$shortcut'";
                shortcut=$basename
                icon="$shortcut"
            }
            $APPS_ARR += $appInfo
        }
    }
}
$content = (ConvertTo-Json $APPS_ARR -Depth 2 | Out-String)
$content | Out-String
