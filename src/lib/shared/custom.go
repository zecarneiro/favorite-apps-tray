package shared

import (
	"main/src/entities"
	"sort"
	"strings"

	"github.com/zecarneiro/golangutils"
)

var (
	APP_NAME      = "Favorites Apps tray"
	ExecutableDir = golangutils.GetExecutableDir()
	ScriptsDir    = golangutils.ResolvePath(golangutils.GetExecutableDir() + "/scripts")
	VendorDir     = golangutils.ResolvePath(golangutils.GetExecutableDir() + "/vendor")
)

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
