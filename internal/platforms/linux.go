//go:build linux

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
)

var (
	linuxAppsInfo []entities.AppsInfo
)

func extractIcon(appInfo entities.AppsInfo) string {
	iconFile := getIcon(appInfo.Shortcut)
	if !file.IsFile(iconFile) {
		if file.IsFile(appInfo.Icon) {
			if file.IsFileExtension(appInfo.Icon, "svg") {
				if !shared.ExtractIcon(appInfo.Icon, iconFile) {
					logger.ErrorStr("Extracted icon: " + iconFile)
				}
			} else if file.IsFileExtension(appInfo.Icon, "png") {
				file.CopyFile(appInfo.Icon, iconFile)
			} else {
				iconFile = ""
			}
		} else {
			iconFile = ""
		}
	}
	return iconFile
}

func getItemInfo(item entities.MenuItemJson) (entities.ItemInfo, error) {
	switch item.Type {
	case enums.SHORTCUTS:
		for _, appInfo := range linuxAppsInfo {
			if appInfo.Shortcut == item.Name+".desktop" || matchRegexByItem(item, appInfo) {
				return getInfoFunc(appInfo, item.Command), nil
			}
		}
	case enums.COMMAND:
		return entities.ItemInfo{Exec: item.Command, Name: item.Name}, nil
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func loadAllLinuxApps(forceLoadApps bool) {
	loadAllApps([]enums.TypeApps{enums.SHORTCUTS}, forceLoadApps)
	linuxAppsInfo, _ = file.ReadJsonFile[[]entities.AppsInfo](getAppsInfoJsonFile(enums.SHORTCUTS))
}

func clearData() {
	linuxAppsInfo = []entities.AppsInfo{}
}

func initApp(forceLoadApps bool) {
	clearData()
	loadAllLinuxApps(forceLoadApps)
}

func runApp(itemInfo entities.ItemInfo) {
	command := models.Command{Cmd: itemInfo.Exec, Verbose: shared.EnableLogs, IsThrow: false, UseShell: true, ShellToUse: golangutilsEnums.Bash, IsAsync: true}
	exe.ExecRealTime(command)
}
