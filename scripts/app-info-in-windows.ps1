param(
    [string] $name,
    [switch] $isStart
)

$ALL_APP = (New-Object -ComObject Shell.Application).NameSpace('shell:AppsFolder').Items()

function printAppInfo($name, $path) {
    $path = "shell:AppsFolder\'$path'"
    $appInfo = @{
        name="$name";
        command="Start-Process $path";
        exec="$path"
    }
    ConvertTo-Json $appInfo -Depth 1
}

function main {
    (New-Object -ComObject Shell.Application).NameSpace('shell:AppsFolder').Items() | ForEach-Object {
        $nameApp = $_.Name
        $pathApp = $_.Path
        if (($isStart -and "$nameApp" -clike "$name*") -or ("$nameApp" -ceq "$name")) {
            printAppInfo "$nameApp" "$pathApp"
            exit 0
        }
    }
}
main