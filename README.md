# favorite-apps-tray

## Run

- Download release package on this [url](https://github.com/zecarneiro/favorite-apps-tray/releases) with name **favorite-apps-tray-VERSION.zip**
- Extract zip file

### Windows

**NOTE:** If you have any problems, please run: `Set-ExecutionPolicy -ExecutionPolicy ByPass -Scope CurrentUser`

1. Open **POWERSHELL** in current directory of this **README**
2. Run executable `favorite-apps-tray.exe`

### Linux

1. Open `terminal` in current directory of this **README**
2. Run `chmod -R 777 .`
3. Run executable `favorite-apps-tray`

**Note:** This application require zenity installed if you use any linux distro.

```bash
# Debian/Ubuntu/Mint
sudo apt install zenity
sudo apt install inkscape
````

## Build local

To build the project its necessary to install **golang** and **git**

```
git clone https://github.com/zecarneiro/favorite-apps-tray.git
```

### Windows

**Note:** If you have any problem with the `nake.ps1` script, please run `powershell -noexit -ExecutionPolicy Bypass -File .\make.ps1 ARGS`.

Run:

```powershell
.\make.ps1 -build
```

### Linux

Run:

```bash
.\make.sh -build
```

## Create Release Packages

Will create a **zip** file on release directory. The zip file contains files on the release directory

### Windows

```powershell
.\make.ps1 -release
```

### Linux

```bash
.\make.sh -release
```
