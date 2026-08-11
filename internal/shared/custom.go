package shared

import (
	"favoriteappstray/internal/entities"
	"fmt"
	"golangutils/pkg/common"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/platform"
	"golangutils/pkg/system"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

func GetExecutableDir() string {
	if platform.IsLinux() {
		return file.JoinPath(system.HomeUserOptDir(), AppName)
	}
	directory, _ := exe.GetExecutableDir()
	return directory
}

func GetScriptsDir() string {
	return file.JoinPath(GetExecutableDir(), "scripts")
}

func loadAppInformations(line string) {
	if strings.HasPrefix(line, "DISPLAY_NAME") {
		_, after, _ := strings.Cut(line, "=")
		AppDisplayName = after
	} else if strings.HasPrefix(line, "VERSION") {
		_, after, _ := strings.Cut(line, "=")
		AppVersion = after
	} else if strings.HasPrefix(line, "RELEASE_DATE") {
		_, after, _ := strings.Cut(line, "=")
		AppReleaseDate = after
	}
}

func setScriptsPermission() {
	if platform.IsWindows() {
		exe.Chmod777(GetExecutableDir(), false)
	} else if platform.IsLinux() {
		exe.Chmod777(GetExecutableDir(), false)
	}
}

func GetConfigIcon(appName string) string {
	iconDir := file.JoinPath(GetConfigurationDir(), "icon")
	file.CreateDirectory(iconDir, true)
	return file.JoinPath(iconDir, appName)
}

func GetLogFile() string {
	return file.JoinPath(GetConfigurationDir(), "log")
}

func GetConfigurationDir() string {
	configDir := file.JoinPath(system.HomeDir(), fmt.Sprintf(".config/%s", AppName))
	file.CreateDirectory(configDir, true)
	return configDir
}

func GetJsonFile() string {
	return file.JoinPath(GetConfigurationDir(), "tray_menu_entries.json")
}

func GetIcon() string {
	icon := ""
	if platform.IsWindows() {
		icon = fmt.Sprintf(`%s/win.ico`, GetExecutableDir())
	} else if platform.IsLinux() {
		icon = fmt.Sprintf(`%s/linux.png`, GetExecutableDir())
	}
	return file.ResolvePath(icon)
}

func SortMenuItemByName(menuItem []entities.MenuItemJson) []entities.MenuItemJson {
	sort.Slice(menuItem, func(i, j int) bool {
		return strings.ToLower(menuItem[i].Name) < strings.ToLower(menuItem[j].Name)
	})
	return menuItem
}

func IsValidateExtension(file string, extensions []string) bool {
	fileBasename := filepath.Base(file)
	fileExtension := filepath.Ext(fileBasename)
	return slices.Contains(extensions, fileExtension)
}

func LoadAppInformations() {
	logger.Error(file.ReadFileLineByLine(file.JoinPath(GetExecutableDir(), "APP_INFO.conf"), loadAppInformations))
	common.WithAppId(AppName)
	logFile := GetLogFile()
	file.DeleteFile(logFile)
	logger.WithLogFile(logFile)
	setScriptsPermission()
}
