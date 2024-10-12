package platform

import (
	"errors"
	"main/src/entities"
	"main/src/lib/golangutils"
	"main/src/lib/shared"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func getIcon(name string) string {
	if len(name) > 0 {
		iconName := filepath.Base(name)
		iconName = strings.TrimSuffix(iconName, filepath.Ext(iconName))
		if golangutils.IsWindows() {
			return shared.GetConfigIcon(iconName + ".ico")
		} else if golangutils.IsLinux() {
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
		shared.ErrorNotify("Your regex: " + item.Regex + ", is faulty")
		return false
	}
	return match
}

func getInfoFunc(app entities.AppsInfo) entities.ItemInfo {
	var icon string
	if golangutils.IsWindows() {
		icon = extractWindowsIcon(app)
	} else if golangutils.IsLinux() {
		icon = extractLinuxIcon(app)
	}
	exec := app.Command
	if golangutils.IsLinux() {
		exec = exec + " &"
	}
	return entities.ItemInfo{Exec: exec, Name: app.DisplayName, Icon: icon}
}

func loadAllApps(cmd string, typeApps []string, force bool) {
	cmdInfo := golangutils.CommandInfo{
		Cmd:     cmd,
		Cwd:     shared.ScriptsDir,
		EnvVars: os.Environ(),
		Verbose: true,
		IsThrow: false,
	}
	if golangutils.IsWindows() {
		cmdInfo.UsePowerShell = true
	} else if golangutils.IsLinux() {
		cmdInfo.UseBash = true
	}
	for _, typeApp := range typeApps {
		// Load shortcuts apps
		if !golangutils.FileExist(getAppsInfoJsonFile(typeApp)) || force {
			cmdInfo.Args = []string{shared.ApplicationName, typeApp}
			golangutils.ExecRealTime(cmdInfo)
		}
	}
}

func getAppsInfoJsonFile(typeApps string) string {
	return golangutils.ResolvePath(shared.GetConfigurationDir() + "/apps-info-" + typeApps + ".json")
}

func GetItemInfo(item entities.MenuItemJson) (entities.ItemInfo, error) {
	if golangutils.IsWindows() {
		return getItemInfoWindows(item)
	} else if golangutils.IsLinux() {
		return getItemInfoLinux(item)
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func Validate() {
	if !golangutils.IsLinux() && !golangutils.IsWindows() {
		shared.ErrorNotify("Invalid Platform.")
		os.Exit(1)
	}
}

func ClearData() {
	clearDataWindows()
	clearDataLinux()
}

func InitPlatform(forceLoadApps bool) {
	if golangutils.IsWindows() {
		initWindows(forceLoadApps)
	} else if golangutils.IsLinux() {
		initLinux(forceLoadApps)
	}
}
