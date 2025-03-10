package platform

import (
	"errors"
	"golangutils"
	"golangutils/entity"
	"main/src/entities"
	"main/src/enums"
	"main/src/lib/shared"
)

var (
	linuxAppsInfo []entities.AppsInfo
)

func extractLinuxIcon(appInfo entities.AppsInfo) string {
	iconFile := getIcon(appInfo.Shortcut)
	if !golangutils.FileExist(iconFile) {
		if golangutils.FileExist(appInfo.Icon) {
			if shared.IsValidateExtension(appInfo.Icon, []string{".svg"}) {
				if !shared.ExtractIcon(appInfo.Icon, iconFile) {
					shared.LoggerUtils.Error("Extracted icon: " + iconFile)
				}
			} else if shared.IsValidateExtension(appInfo.Icon, []string{".png"}) {
				golangutils.CopyFile(appInfo.Icon, iconFile)
			} else {
				iconFile = ""
			}
		} else {
			iconFile = ""
		}
	}
	return iconFile
}

func getItemInfoLinux(item entities.MenuItemJson) (entities.ItemInfo, error) {
	if item.Type == enums.SHORTCUTS {
		for _, appInfo := range linuxAppsInfo {
			if appInfo.Shortcut == item.Name+".desktop" || matchRegexByItem(item, appInfo) {
				return getInfoFunc(appInfo, item.Command), nil
			}
		}
	} else if item.Type == enums.COMMAND {
		return entities.ItemInfo{Exec: item.Command, Name: item.Name}, nil
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func loadAllLinuxApps(forceLoadApps bool) {
	loadAllApps([]string{enums.SHORTCUTS}, forceLoadApps)
	linuxAppsInfo, _ = golangutils.ReadJsonFile[[]entities.AppsInfo](getAppsInfoJsonFile(enums.SHORTCUTS))
}

func clearDataLinux() {
	linuxAppsInfo = []entities.AppsInfo{}
}

func initLinux(forceLoadApps bool) {
	clearDataLinux()
	loadAllLinuxApps(forceLoadApps)
}

func runAppLinux(itemInfo entities.ItemInfo) {
	command := entity.Command{Cmd: itemInfo.Exec, Verbose: shared.EnableLogs, IsThrow: false, UseBash: true}
	shared.ConsoleUtils.ExecAsync(command, shared.ProcessConsoleResult)
}
