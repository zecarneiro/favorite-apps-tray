package platform

import (
	"main/src/entities"
	"main/src/lib/shared"
	"os"
	"path/filepath"
	"strings"

	"github.com/zecarneiro/golangutils"
)

func getIcon(info entities.AppsInfoWindows) string {
	iconName := filepath.Base(info.Name)
	return shared.GetConfigIcon(strings.TrimSuffix(iconName, filepath.Ext(iconName)) + ".ico")
}

func GetItemInfo(item entities.MenuItemJson) (entities.ItemInfo, error) {
	if golangutils.IsWindows() {
		return getItemInfoWindows(item)
	}
	return entities.ItemInfo{}, nil
}

func Validate() {
	if !golangutils.IsLinux() && !golangutils.IsWindows() {
		shared.ErrorNotify("Invalid Platform.")
		os.Exit(1)
	}
}

func InitPlatform() {
	if golangutils.IsWindows() {
		initWindows()
	} else if golangutils.IsLinux() {
		shared.InfoNotify("Not implemented yet")
		os.Exit(0)
	}
}
