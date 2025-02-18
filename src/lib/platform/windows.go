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
	windowsAppsInfo   []entities.AppsInfo
	shortcutsAppsInfo []entities.AppsInfo
)

func extractWindowsIcon(appInfo entities.AppsInfo) string {
	iconFile := getIcon(appInfo.Shortcut)
	if !golangutils.FileExist(iconFile) {
		commandInfo := entity.Command{
			Cmd:           ". \"" + golangutils.ResolvePath(shared.VendorDir+"/powershell-utils/MainUtils.ps1") + "\"",
			IsThrow:       false,
			UsePowerShell: true,
			Args:          []string{"IMPORT_ALL_LIBS; icon_extractor -file \"" + appInfo.Icon + "\" -dest \"" + iconFile + "\" -display"},
			Verbose:       shared.EnableLogs,
		}
		shared.ConsoleUtils.ExecRealTime(commandInfo)
	}
	return iconFile
}

func getItemInfoWindows(item entities.MenuItemJson) (entities.ItemInfo, error) {
	if item.Type == enums.WINDOWS_APPS {
		for _, appInfo := range windowsAppsInfo {
			if appInfo.DisplayName == item.Name || matchRegexByItem(item, appInfo) {
				return getInfoFunc(appInfo, item.Command), nil
			}
		}
	} else if item.Type == enums.SHORTCUTS {
		for _, appInfo := range shortcutsAppsInfo {
			if appInfo.Shortcut == item.Name+".lnk" || matchRegexByItem(item, appInfo) {
				return getInfoFunc(appInfo, item.Command), nil
			}
		}
	} else if item.Type == enums.COMMAND {
		return entities.ItemInfo{Exec: item.Command, Name: item.Name}, nil
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func loadAllWindowsApps(forceLoadApps bool) {
	loadAllApps(golangutils.ResolvePath(shared.ScriptsDir+"/app-info-in-windows.ps1"), []string{enums.SHORTCUTS, enums.WINDOWS_APPS}, forceLoadApps)
	shortcutsAppsInfo, _ = golangutils.ReadJsonFile[[]entities.AppsInfo](getAppsInfoJsonFile(enums.SHORTCUTS))
	windowsAppsInfo, _ = golangutils.ReadJsonFile[[]entities.AppsInfo](getAppsInfoJsonFile(enums.WINDOWS_APPS))
}

func clearDataWindows() {
	windowsAppsInfo = []entities.AppsInfo{}
	shortcutsAppsInfo = []entities.AppsInfo{}
}

func initWindows(forceLoadApps bool) {
	clearDataWindows()
	loadAllWindowsApps(forceLoadApps)
}
