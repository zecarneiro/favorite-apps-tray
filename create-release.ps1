# Author: José M. C. Noronha

function remove_directories($directory) {
    if (Test-Path -Path "$directory") {
        Write-Host "Remove directory: '$directory'"
        Remove-Item "$directory" -Recurse -Force
    }
}
remove_directories "$pwd\node_modules"
remove_directories "$pwd\dist"
remove_directories "$pwd\out"

Write-Host "Install dependencies..."
npm install

if ($?) {
    Write-Host "Build Project..."
    npm run build
}
if ($?) {
    Write-Host "Create release..."
    npm run make
}
Write-Host "Done."