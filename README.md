# favorite-apps-tray

---

**NOTE**

If is the first time that you run this application, you should follow the instructions:

- Windows
  1. Open _CMD_ in current directory of this **README**
  2. Run `scripts\set-powershell-exec-policy.bat`
- Linux
  1. Open `terminal` in current directory of this **README**
  2. Run `chmod -R 777 .`

---

## Instalation

### Windows

```
.\make.ps1 -install
```

To uninstall run `.\make.ps1 -uninstall`

### Linux

```
.\make.sh -install
```

To uninstall run `.\make.sh -uninstall`

## Build local

To build the project its necessary to install **golang** and **git**

```
git clone https://github.com/zecarneiro/favorite-apps-tray.git
```

### Windows

Run:

```
.\make.ps1 -installDependencies
.\make.ps1 -build
```

### Ubuntu

Run:

```
.\make.sh -installDependencies
.\make.sh -build
```

## Create Release Packages

Will create a **zip** file on release directory. The zip file contains files on the release directory

### Windows

```
.\make.ps1 -release
```

### Ubuntu

```
.\make.sh -release
```