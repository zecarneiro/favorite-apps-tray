package platform

import (
	"errors"
	"main/src/entities"
	"main/src/enums"
	"main/src/lib/shared"
	"path/filepath"
	"strings"

	"github.com/zecarneiro/golangutils"
)

var (
	windowsAppsInfo   = []entities.AppsInfoWindows{}
	shortcutsAppsInfo = []entities.AppsInfoWindows{}
)

func extractWindowsIcon(appInfo entities.AppsInfoWindows) {
	iconFile := getIcon(appInfo)
	if !golangutils.FileExist(iconFile) {
		powershellUtilsScript := golangutils.ResolvePath(shared.VendorDir + "/powershell-utils/MainUtils.ps1")
		commandInfo := golangutils.CommandInfo{
			Cmd:           ". \"" + powershellUtilsScript + "\"",
			IsThrow:       false,
			UsePowerShell: true,
			Args:          []string{";icon_extractor -file \"" + appInfo.Exec + "\" -dest \"" + iconFile + "\" -display"},
			Verbose:       false,
		}
		golangutils.ExecRealTime(commandInfo)
	}
}

func getItemInfoWindows(item entities.MenuItemJson) (entities.ItemInfo, error) {
	if item.Type == enums.WINDOWS_APPS {
		for _, appInfo := range windowsAppsInfo {
			if (item.IsInitialName && strings.HasPrefix(appInfo.Name, item.Name)) || appInfo.Name == item.Name {
				extractWindowsIcon(appInfo)
				return entities.ItemInfo{Exec: appInfo.Command, Name: appInfo.Name, Icon: getIcon(appInfo)}, nil
			}
		}
	} else if item.Type == enums.SHORTCUTS {
		for _, appInfo := range shortcutsAppsInfo {
			if (item.IsInitialName && strings.HasPrefix(filepath.Base(appInfo.Name), item.Name)) || filepath.Base(appInfo.Name) == item.Name+".lnk" {
				extractWindowsIcon(appInfo)
				return entities.ItemInfo{Exec: appInfo.Command, Name: item.Name, Icon: getIcon(appInfo)}, nil
			}
		}
	} else if item.Type == enums.COMMAND {
		return entities.ItemInfo{Exec: item.Command, Name: item.Name}, nil
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func loadAllShortcuts() {
	startMenuDirs := []string{
		golangutils.ResolvePath(golangutils.SysInfo().HomeDir + "/AppData/Roaming/Microsoft/Windows/Start Menu"),
		"C:\\ProgramData\\Microsoft\\Windows\\Start Menu\\Programs",
	}
	for _, menuDir := range startMenuDirs {
		res, err := golangutils.ReadDirRecursive(menuDir)
		if err != nil {
			shared.ErrorNotify(err.Error())
		} else {
			for _, appInfo := range res.Files {
				shortcutsAppsInfo = append(shortcutsAppsInfo, entities.AppsInfoWindows{Exec: appInfo, Command: "Start-Process -FilePath '" + appInfo + "'", Name: appInfo})
			}
		}
	}
}

func loadAllWindowsApps() {
	res := golangutils.Exec(golangutils.CommandInfo{Cmd: golangutils.ResolvePath(shared.ScriptsDir + "/app-info-in-windows.ps1"), UsePowerShell: true, IsThrow: false})
	if res.HasError() {
		shared.ErrorNotify(res.Error.Error())
	} else {
		windowsAppsInfo, _ = golangutils.StringToObject[[]entities.AppsInfoWindows](res.Data)
	}
}

func initWindows() {
	loadAllWindowsApps()
	loadAllShortcuts()
}
