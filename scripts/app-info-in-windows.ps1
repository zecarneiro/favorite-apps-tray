$ALL_APP = (New-Object -ComObject Shell.Application).NameSpace('shell:AppsFolder').Items()

function main {
    $appsArr = @()
    $ALL_APP | ForEach-Object {
        $nameApp = $_.Name
        $pathApp = $_.Path
        $path = "shell:AppsFolder\'$pathApp'"
        $appInfo = @{
            name="$nameApp";
            command="Start-Process -FilePath $path";
            exec="$path"
        }
        $appsArr += $appInfo
    }
    ConvertTo-Json $appsArr -Depth 1
}
main
