package lib

import (
	"main/src/entities"
	"main/src/enums"
	"main/src/lib/platform"
	"main/src/lib/shared"
	"sort"

	"github.com/energye/systray"
	"github.com/sqweek/dialog"
	"github.com/zecarneiro/golangutils"
)

var (
	menuJsonData entities.MenuJson
	menu         []*systray.MenuItem
)

func refresh() {
	systray.ResetMenu()
	processApp()
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
	isItemInserted := false
	for _, item := range items {
		canProcessMenuItem := false
		if isValidItem(item) {
			menuItem := &systray.MenuItem{}
			appInfo, err := platform.GetItemInfo(item)
			if err == nil {
				if mainMenu != nil {
					isItemInserted = true
					canProcessMenuItem = true
					menuItem = mainMenu.AddSubMenuItem(item.Name, item.Name)
				} else {
					isItemInserted = true
					canProcessMenuItem = true
					menuItem = systray.AddMenuItem(item.Name, item.Name)
				}
				if canProcessMenuItem {
					if len(appInfo.Icon) > 0 && golangutils.FileExist(appInfo.Icon) {
						menuItem.SetIcon(golangutils.ReadFileInByte(appInfo.Icon))
					}
					menuItem.Click(func() {
						command := golangutils.CommandInfo{Cmd: appInfo.Exec, Verbose: false, IsThrow: false}
						if golangutils.IsWindows() {
							command.UsePowerShell = true
						}
						golangutils.ExecRealTimeAsync(command)
					})
				}
			}
		}
	}
	if !isItemInserted {
		if mainMenu != nil {
			mainMenu.Disable()
		} else {
			buildEmptyMenu()
		}
	}
}

func buildSettingMenu() {
	settingsMenu := systray.AddMenuItem("Settings", "Settings")
	settingsMenu.AddSubMenuItem("Update Menu", "Update Menu for any changes").Click(func() {
		refresh()
	})
	settingsMenu.AddSubMenuItem("Select/Change JSON file", "Select JSON configuration file").Click(func() {
		filename, err := dialog.File().Title("Select JSON configuration file").Load()
		if err != nil {
			shared.ErrorNotify(err.Error())
		} else {
			golangutils.DeleteFile(shared.GetJsonFile())
			err := golangutils.CopyFile(filename, shared.GetJsonFile())
			if err != nil {
				shared.ErrorNotify(err.Error())
			} else {
				refresh()
			}
		}
	})
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
	systray.AddSeparator()
	buildMenuItem(menuJsonData.NoMenu, nil)
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

func processApp() {
	loadMenuJsonData()
	platform.InitPlatform()
	buildMenu()
}

func buildTrayApp() {
	systray.SetIcon(golangutils.ReadFileInByte(shared.GetIcon()))
	systray.SetTitle(shared.APP_NAME)
	systray.SetTooltip(shared.APP_NAME)
	systray.SetOnClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})
	systray.SetOnRClick(func(menu systray.IMenu) {
		menu.ShowMenu()
	})
}

func Start() {
	platform.Validate()
	buildTrayApp()
	processApp()
}
