package internal

import (
	"favoriteappstray/internal/entities"
	"favoriteappstray/internal/enums"
	"favoriteappstray/internal/platforms"
	"favoriteappstray/internal/shared"
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/ui"
	"sort"

	"github.com/energye/systray"
)

var (
	menuJsonData     entities.MenuJson
	menu             []*systray.MenuItem
	isSystrayCreated = false
)

func refresh(forceLoadApps bool) {
	platforms.InitPlatform(forceLoadApps)
	systray.ResetMenu()
	loadMenuJsonData()
	buildTrayApp()
	updateMenuJsonData()
}

func isValidItem(item entities.MenuItemJson) bool {
	if len(item.Name) < 1 {
		shared.ErrorNofity("Invalid Item: " + item.Name)
		return false
	}
	if item.Type != enums.WINDOWS_APPS && item.Type != enums.SHORTCUTS && item.Type != enums.COMMAND {
		shared.ErrorNofity(fmt.Sprintf("Invalid Item Type: %s, from name: %s", item.Type, item.Name))
		return false
	}
	return true
}

func buildMenuItem(items []entities.MenuItemJson, mainMenu *systray.MenuItem) {
	itemsInfo := []entities.ItemInfo{}
	for _, item := range items {
		if isValidItem(item) {
			appInfo, err := platforms.GetItemInfo(item)
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
			if len(itemInfo.Icon) > 0 && file.IsFile(itemInfo.Icon) {
				icon, err := file.ReadFileInByte(itemInfo.Icon)
				if err != nil {
					shared.ErrorNofity(err.Error())
				} else {
					menuItem.SetIcon(icon)
				}
			}
			menuItem.Click(func() {
				platforms.RunApp(itemInfo)
			})
		}
	}
}

func buildSettingMenu() {
	settingsMenu := systray.AddMenuItem("Settings", "Settings")
	settingsMenu.AddSubMenuItem("Update Menu", "Update Menu for any changes").Click(func() {
		shared.ShowProcessingMsg(false)
		refresh(true)
		shared.ShowProcessingMsg(true)
	})
	settingsMenu.AddSubMenuItem("Select/Change JSON file", "Select JSON configuration file").Click(func() {
		filenameResp := ui.SelectFile(shared.AppName)
		if filenameResp.HasError() {
			shared.ErrorNofity(filenameResp.Error.Error())
		} else {
			shared.ShowProcessingMsg(false)
			file.DeleteFile(shared.GetJsonFile())
			if err := file.CopyFile(filenameResp.Data, shared.GetJsonFile()); err != nil {
				shared.ErrorNofity(err.Error())
			} else {
				refresh(true)
			}
			shared.ShowProcessingMsg(true)
		}
	})
	enableLogsItem := settingsMenu.AddSubMenuItemCheckbox("Enable Logs", "Enable logs for most of operations", menuJsonData.EnableLogs)
	enableLogsItem.Click(func() {
		message := ""
		if enableLogsItem.Checked() {
			enableLogsItem.Uncheck()
			message = "disabled"
			shared.EnableLogs = false
		} else {
			enableLogsItem.Check()
			message = "enabled"
			shared.EnableLogs = true
		}
		ui.WithVerbose(shared.EnableLogs)
		menuJsonData.EnableLogs = shared.EnableLogs
		updateMenuJsonData()
		shared.InfoNofity(fmt.Sprintf("All Logs was %s by user.", message))
	})
	// About Settings
	aboutSettings := settingsMenu.AddSubMenuItem("About", "About")
	aboutSettings.Click(func() {
		message := "Name: " + shared.AppName
		message += "\nVersion: " + shared.AppVersion
		message += "\nRelease Date: " + shared.AppReleaseDate
		message += "\nLog file located: " + shared.GetLogFile()
		shared.InfoDialog(message)
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
	if len(menuJsonData.NoMenu) > 0 {
		systray.AddSeparator()
		buildMenuItem(menuJsonData.NoMenu, nil)
	}
	platforms.ClearData()
	systray.AddSeparator()
	buildSettingMenu()
	systray.AddMenuItem("Exit", "Exit of the application").Click(func() {
		systray.Quit()
	})
}

func loadMenuJsonData() {
	menuJsonData = entities.MenuJson{}
	if file.IsFile(shared.GetJsonFile()) {
		data, err := file.ReadJsonFile[entities.MenuJson](shared.GetJsonFile())
		if err != nil {
			shared.ErrorNofity(err.Error())
		} else {
			ui.WithVerbose(shared.EnableLogs)
			menuJsonData = data
			menuJsonData.EnableLogs = shared.EnableLogs
			menuJsonData.NoMenu = shared.SortMenuItemByName(menuJsonData.NoMenu)
			for _, othersMenuItem := range menuJsonData.Others {
				othersMenuItem = shared.SortMenuItemByName(othersMenuItem)
			}
		}
	}
}

func updateMenuJsonData() {
	if err := file.WriteJsonFile(shared.GetJsonFile(), menuJsonData, false); err != nil {
		shared.ErrorNofity(err.Error())
	}
}

func showMenu(menu systray.IMenu) {
	if menu != nil {
		menu.ShowMenu()
	} else {
		shared.ShowProcessingMsg(false)
	}
}

func buildTrayApp() {
	buildMenu()
	if !isSystrayCreated {
		fileByte, _ := file.ReadFileInByte(shared.GetIcon())
		systray.SetIcon(fileByte)
		systray.SetTitle(shared.AppName)
		systray.SetTooltip(shared.AppName)
		systray.SetOnClick(showMenu)
		systray.SetOnRClick(showMenu)
		isSystrayCreated = true
		refresh(false)
	}
}

func Start() {
	platforms.Validate()
	shared.LoadAppInformations()
	ui.WithVerbose(shared.EnableLogs)
	platforms.InitPlatform(false)
	systray.Run(buildTrayApp, nil)
}
