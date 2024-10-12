# Favorite Apps Tray

- This repository contains app to add an menu with favorites apps on system tray

## Configurations

### Windows
- All **Shurtcuts** , its on:
  - On windows
    - `C:\Users\YOUR_USER_NAME\AppData\Roaming\Microsoft\Windows\Start Menu\Programs`
    - `C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs`
  - On Linux
    ```
    ~/.local/share/applications
    /usr/share/applications
    /var/lib/snapd/desktop/applications
    /var/lib/flatpak/exports/share/applications
    ~/.local/share/flatpak/exports/share/applications
    ```

**NOTE**
- If **Shurtcuts** not exist, the content not appear


### For Electron Forge
Run:

- `npm install @electron-forge/cli @electron-forge/maker-deb @electron-forge/maker-rpm @electron-forge/maker-squirrel @electron-forge/maker-zip @electron-forge/shared-types -D`
- `npx electron-forge import`
- Create `forge.config.js` on your root project directory and insert
```javascript
...
module.exports = {
  "packagerConfig": {
    "icon": "./assets/image/icon"
  },
  "makers": [
    {
      "name": "@electron-forge/maker-squirrel",
      "config": {}
    },
    {
      "name": "@electron-forge/maker-zip",
      "platforms": ["darwin"]
    },
    {
      "name": "@electron-forge/maker-deb",
      "config": {}
    },
    {
      "name": "@electron-forge/maker-rpm",
      "config": {}
    }
  ]
}
```
