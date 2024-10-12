package shared

import (
	"main/src/entities"
	"main/src/lib/golangutils"
	"path/filepath"
	"sort"
	"strings"
)

var (
	ExecutableDir          = golangutils.GetExecutableDir()
	ScriptsDir             = golangutils.ResolvePath(golangutils.GetExecutableDir() + "/scripts")
	VendorDir              = golangutils.ResolvePath(golangutils.GetExecutableDir() + "/vendor")
	ApplicationName        string
	ApplicationVersion     string
	ApplicationReleaseDate string
)

func loadAppInformations(line string, err error) {
	golangutils.ProcessError(err)
	if strings.HasPrefix(line, "NAME") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationName = after
	} else if strings.HasPrefix(line, "VERSION") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationVersion = after
	} else if strings.HasPrefix(line, "RELEASE_DATE") {
		_, after, _ := strings.Cut(line, "=")
		ApplicationReleaseDate = after
	}
}

func GetConfigIcon(appName string) string {
	iconDir := golangutils.ResolvePath(GetConfigurationDir() + "/icon")
	golangutils.CreateDirectory(iconDir, true)
	return golangutils.ResolvePath(iconDir + "/" + appName)
}

func GetConfigurationDir() string {
	configDir := golangutils.ResolvePath(golangutils.SysInfo().HomeDir + "/.config/favorite-apps-tray")
	golangutils.CreateDirectory(configDir, true)
	return configDir
}

func GetJsonFile() string {
	return golangutils.ResolvePath(GetConfigurationDir() + "/tray_menu_entries.json")
}

func GetIcon() string {
	icon := ExecutableDir + "/assets/image/logo/icon"
	if golangutils.IsWindows() {
		icon = icon + ".ico"
	} else if golangutils.IsLinux() {
		icon = icon + ".png"
	}
	return golangutils.ResolvePath(icon)
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
	return golangutils.InArray(extensions, fileExtension)
}

func LoadAppInformations() {
	golangutils.ReadFileLineByLine(golangutils.ResolvePath(ExecutableDir+"/APP_INFO.conf"), loadAppInformations)
}
