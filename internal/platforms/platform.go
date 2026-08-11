package platforms

import (
	"errors"
	"favoriteappstray/internal/entities"
	"favoriteappstray/internal/enums"
	"favoriteappstray/internal/shared"
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/logic"
	"golangutils/pkg/platform"
	"path/filepath"
	"regexp"
	"strings"
)

func getIcon(name string) string {
	if len(name) > 0 {
		iconName := filepath.Base(name)
		iconName = strings.TrimSuffix(iconName, filepath.Ext(iconName))
		if platform.IsWindows() {
			return shared.GetConfigIcon(iconName + ".ico")
		} else if platform.IsLinux() {
			return shared.GetConfigIcon(iconName + ".png")
		}
	}
	return ""
}

func matchRegexByItem(item entities.MenuItemJson, appInfo entities.AppsInfo) bool {
	var match bool
	var err error
	if len(item.Regex) == 0 {
		return false
	}
	if item.RegexOnDisplayName {
		match, err = regexp.MatchString(item.Regex, appInfo.DisplayName)
	} else {
		match, err = regexp.MatchString(item.Regex, appInfo.Shortcut)
	}
	if err != nil {
		shared.ErrorNofity("Your regex: " + item.Regex + ", is faulty")
		return false
	}
	return match
}

func getInfoFunc(app entities.AppsInfo, defaultCommand string) entities.ItemInfo {
	var icon string
	exec := app.Command
	if platform.IsWindows() || platform.IsLinux() {
		icon = extractIcon(app)
	}
	if len(defaultCommand) > 0 {
		exec = defaultCommand
	}
	if platform.IsLinux() {
		exec = exec + " &"
	}
	return entities.ItemInfo{Exec: exec, Name: app.DisplayName, Icon: icon}
}

func loadAllApps(typeApps []enums.TypeApps, force bool) {
	for _, typeApp := range typeApps {
		// Load shortcuts apps
		if !file.FileExist(getAppsInfoJsonFile(typeApp)) || force {
			shared.AppsInfo(typeApp)
		}
	}
}

func getAppsInfoJsonFile(typeApps enums.TypeApps) string {
	return file.JoinPath(shared.GetConfigurationDir(), fmt.Sprintf("apps-info-%s.json", typeApps))
}

func GetItemInfo(item entities.MenuItemJson) (entities.ItemInfo, error) {
	if platform.IsWindows() || platform.IsLinux() {
		return getItemInfo(item)
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func Validate() {
	if !platform.IsLinux() && !platform.IsWindows() {
		shared.ErrorNofity("Invalid Platform.")
		logic.Exit(1)
	}
}

func ClearData() {
	if platform.IsWindows() || platform.IsLinux() {
		clearData()
	}
}

func InitPlatform(forceLoadApps bool) {
	if platform.IsWindows() || platform.IsLinux() {
		initApp(forceLoadApps)
	}
}

func RunApp(itemInfo entities.ItemInfo) {
	if platform.IsWindows() || platform.IsLinux() {
		runApp(itemInfo)
	}
}
