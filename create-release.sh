#!/bin/bash
# Author: José M. C. Noronha

remove_directories() {
    local directory="$1"
    if [ -d "$directory" ]; then
        echo "Remove directory: '$directory'"
        rm -rf "$directory"
    fi
}
remove_directories "$PWD/node_modules"
remove_directories "$PWD/dist"
remove_directories "$PWD/out"

install_dependencies() {
    echo "Install dependencies..."
    npm install
}
build_project() {
    echo "Build Project..."
    npm run build
}
create_release() {
    echo "Create release..."
    npm run make
}

install_dependencies && build_project && create_release
echo "Done."