//go:build windows

package platforms

import (
	"errors"
	"favoriteappstray/internal/entities"
	"favoriteappstray/internal/enums"
	"favoriteappstray/internal/shared"
	golangutilsEnums "golangutils/pkg/enums"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/models"
	"path/filepath"
)

var (
	windowsAppsInfo   []entities.AppsInfo
	shortcutsAppsInfo []entities.AppsInfo
)

func extractIcon(appInfo entities.AppsInfo) string {
	iconFile := getIcon(appInfo.Shortcut)
	if !file.IsFile(iconFile) {
		fileType := filepath.Ext(appInfo.Icon)
		if fileType == ".lnk" || fileType == ".exe" {
			if !shared.ExtractIcon(appInfo.Icon, iconFile) {
				logger.ErrorStr("Extracted icon: " + iconFile)
			}
		} else {
			iconFile = ""
		}
	}
	return iconFile
}

func getItemInfo(item entities.MenuItemJson) (entities.ItemInfo, error) {
	switch item.Type {
	case enums.WINDOWS_APPS:
		for _, appInfo := range windowsAppsInfo {
			if appInfo.DisplayName == item.Name || matchRegexByItem(item, appInfo) {
				return getInfoFunc(appInfo, item.Command), nil
			}
		}
	case enums.SHORTCUTS:
		for _, appInfo := range shortcutsAppsInfo {
			if appInfo.Shortcut == item.Name+".lnk" || matchRegexByItem(item, appInfo) {
				return getInfoFunc(appInfo, item.Command), nil
			}
		}
	case enums.COMMAND:
		return entities.ItemInfo{Exec: item.Command, Name: item.Name}, nil
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func loadAllWindowsApps(forceLoadApps bool) {
	loadAllApps([]enums.TypeApps{enums.SHORTCUTS, enums.WINDOWS_APPS}, forceLoadApps)
	shortcutsAppsInfo, _ = file.ReadJsonFile[[]entities.AppsInfo](getAppsInfoJsonFile(enums.SHORTCUTS))
	windowsAppsInfo, _ = file.ReadJsonFile[[]entities.AppsInfo](getAppsInfoJsonFile(enums.WINDOWS_APPS))
}

func clearData() {
	windowsAppsInfo = []entities.AppsInfo{}
	shortcutsAppsInfo = []entities.AppsInfo{}
}

func initApp(forceLoadApps bool) {
	clearData()
	loadAllWindowsApps(forceLoadApps)
}

func runApp(itemInfo entities.ItemInfo) {
	command := models.Command{Cmd: itemInfo.Exec, Verbose: shared.EnableLogs, IsThrow: false, UseShell: true, ShellToUse: golangutilsEnums.PowerShell, IsAsync: true}
	exe.ExecRealTime(command)
}
