package platform

import (
	"errors"
	"golangutils"
	"golangutils/entity"
	"main/src/entities"
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
		if shared.SystemUtils.IsWindows() {
			return shared.GetConfigIcon(iconName + ".ico")
		} else if shared.SystemUtils.IsLinux() {
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

func getInfoFunc(app entities.AppsInfo, defaultCommand string) entities.ItemInfo {
	var icon string
	if shared.SystemUtils.IsWindows() {
		icon = extractWindowsIcon(app)
	} else if shared.SystemUtils.IsLinux() {
		icon = extractLinuxIcon(app)
	}
	exec := app.Command
	if len(defaultCommand) > 0 {
		exec = defaultCommand
	}
	if shared.SystemUtils.IsLinux() {
		exec = exec + " &"
	}
	return entities.ItemInfo{Exec: exec, Name: app.DisplayName, Icon: icon}
}

func loadAllApps(cmd string, typeApps []string, force bool) {
	cmdInfo := entity.Command{
		Cmd:     cmd,
		Cwd:     shared.ScriptsDir,
		EnvVars: os.Environ(),
		Verbose: shared.EnableLogs,
		IsThrow: false,
	}
	if shared.SystemUtils.IsWindows() {
		cmdInfo.UsePowerShell = true
	} else if shared.SystemUtils.IsLinux() {
		cmdInfo.UseBash = true
	}
	for _, typeApp := range typeApps {
		// Load shortcuts apps
		if !golangutils.FileExist(getAppsInfoJsonFile(typeApp)) || force {
			cmdInfo.Args = []string{shared.ApplicationName, typeApp}
			shared.ConsoleUtils.ExecRealTime(cmdInfo)
		}
	}
}

func getAppsInfoJsonFile(typeApps string) string {
	return golangutils.ResolvePath(shared.GetConfigurationDir() + "/apps-info-" + typeApps + ".json")
}

func GetItemInfo(item entities.MenuItemJson) (entities.ItemInfo, error) {
	if shared.SystemUtils.IsWindows() {
		return getItemInfoWindows(item)
	} else if shared.SystemUtils.IsLinux() {
		return getItemInfoLinux(item)
	}
	return entities.ItemInfo{}, errors.New("Not found item info: " + item.Name)
}

func Validate() {
	if !shared.SystemUtils.IsLinux() && !shared.SystemUtils.IsWindows() {
		shared.ErrorNotify("Invalid Platform.")
		os.Exit(1)
	}
}

func ClearData() {
	clearDataWindows()
	clearDataLinux()
}

func InitPlatform(forceLoadApps bool) {
	if shared.SystemUtils.IsWindows() {
		initWindows(forceLoadApps)
	} else if shared.SystemUtils.IsLinux() {
		initLinux(forceLoadApps)
	}
}
