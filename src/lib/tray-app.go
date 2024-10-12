package lib

import (
	"main/src/entities"
	"main/src/enums"
	"main/src/lib/golangutils"
	"main/src/lib/platform"
	"main/src/lib/shared"
	"sort"

	"github.com/energye/systray"
)

var (
	menuJsonData     entities.MenuJson
	menu             []*systray.MenuItem
	isSystrayCreated = false
)

func refresh(forceLoadApps bool) {
	platform.InitPlatform(forceLoadApps)
	systray.ResetMenu()
	loadMenuJsonData()
	buildTrayApp()
}

func isValidItem(item entities.MenuItemJson) bool {
	if len(item.Name) < 1 {
		shared.ErrorNotify("Invalid Item: " + item.Name)
		return false
	}
	if item.Type != enums.WINDOWS_APPS && item.Type != enums.SHORTCUTS && item.Type != enums.COMMAND {
		shared.ErrorNotify("Invalid Item Type: " + item.Type + ", from name: " + item.Name)
		return false
	}
	return true
}

func buildMenuItem(items []entities.MenuItemJson, mainMenu *systray.MenuItem) {
	itemsInfo := []entities.ItemInfo{}
	for _, item := range items {
		if isValidItem(item) {
			appInfo, err := platform.GetItemInfo(item)
			if err == nil {
				itemsInfo = append(itemsInfo, appInfo)
			}
		}
	}
	if len(itemsInfo) == 0 {
		if mainMenu != nil {
			mainMenu.Disable()
		} else {
			buildEmptyMenu()
		}
	} else {
		for _, itemInfo := range itemsInfo {
			var menuItem *systray.MenuItem
			if mainMenu != nil {
				menuItem = mainMenu.AddSubMenuItem(itemInfo.Name, itemInfo.Name)
			} else {
				menuItem = systray.AddMenuItem(itemInfo.Name, itemInfo.Name)
			}
			if len(itemInfo.Icon) > 0 && golangutils.FileExist(itemInfo.Icon) {
				menuItem.SetIcon(golangutils.ReadFileInByte(itemInfo.Icon))
			}
			menuItem.Click(func() {
				command := golangutils.CommandInfo{Cmd: itemInfo.Exec, Verbose: true, IsThrow: false}
				if golangutils.IsWindows() {
					command.UsePowerShell = true
				} else if golangutils.IsLinux() {
					command.UseBash = true
				}
				golangutils.ExecRealTimeAsync(command)
			})
		}
	}
}

func buildSettingMenu() {
	settingsMenu := systray.AddMenuItem("Settings", "Settings")
	settingsMenu.AddSubMenuItem("Update Menu", "Update Menu for any changes").Click(func() {
		shared.InfoNotify("Processing...")
		refresh(true)
		shared.InfoNotify("Processing, done.")
	})
	settingsMenu.AddSubMenuItem("Select/Change JSON file", "Select JSON configuration file").Click(func() {
		filename, err := shared.SelectFileDialog()
		if err != nil {
			shared.ErrorNotify(err.Error())
		} else {
			shared.InfoNotify("Processing...")
			golangutils.DeleteFile(shared.GetJsonFile())
			err := golangutils.CopyFile(filename, shared.GetJsonFile())
			if err != nil {
				shared.ErrorNotify(err.Error())
			} else {
				refresh(true)
			}
			shared.InfoNotify("Processing, done.")
		}
	})

	// About Settings
	aboutSettings := settingsMenu.AddSubMenuItem("About", "About")
	aboutSettings.AddSubMenuItem("Name: "+shared.ApplicationName, "").Disable()
	aboutSettings.AddSubMenuItem("Version: "+shared.ApplicationVersion, "").Disable()
	aboutSettings.AddSubMenuItem("Release Date: "+shared.ApplicationReleaseDate, "").Disable()
}

func buildEmptyMenu() {
	systray.AddMenuItem("Empty", "").Disable()
}

func buildOthersMenu() {
	if len(menuJsonData.Others) == 0 {
		buildEmptyMenu()
	} else {
		keys := make([]string, 0, len(menuJsonData.Others))
		for k := range menuJsonData.Others {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			mainMenu := systray.AddMenuItem(key, key)
			buildMenuItem(menuJsonData.Others[key], mainMenu)
		}
	}
}

func buildMenu() {
	buildOthersMenu()
	if len(menuJsonData.NoMenu) > 0 {
		systray.AddSeparator()
		buildMenuItem(menuJsonData.NoMenu, nil)
	}
	platform.ClearData()
	systray.AddSeparator()
	buildSettingMenu()
	systray.AddMenuItem("Exit", "Exit of the application").Click(func() {
		systray.Quit()
	})
}

func loadMenuJsonData() {
	menuJsonData = entities.MenuJson{}
	if golangutils.FileExist(shared.GetJsonFile()) {
		data, err := golangutils.ReadJsonFile[entities.MenuJson](shared.GetJsonFile())
		if err != nil {
			shared.ErrorNotify(err.Error())
		} else {
			menuJsonData = data
			menuJsonData.NoMenu = shared.SortMenuItemByName(menuJsonData.NoMenu)
			for _, othersMenuItem := range menuJsonData.Others {
				othersMenuItem = shared.SortMenuItemByName(othersMenuItem)
			}
		}
	}
}

func buildTrayApp() {
	buildMenu()
	if !isSystrayCreated {
		systray.SetIcon(golangutils.ReadFileInByte(shared.GetIcon()))
		systray.SetTitle(shared.ApplicationName)
		systray.SetTooltip(shared.ApplicationName)
		systray.SetOnClick(func(menu systray.IMenu) {
			menu.ShowMenu()
		})
		systray.SetOnRClick(func(menu systray.IMenu) {
			menu.ShowMenu()
		})
		isSystrayCreated = true
		refresh(false)
	}
}

func Start() {
	platform.Validate()
	shared.LoadAppInformations()
	platform.InitPlatform(false)
	systray.Run(buildTrayApp, nil)
}
