//go:build windows

package internal

import (
	"syscall"

	"github.com/energye/systray"
)

var (
	modKernel32     = syscall.NewLazyDLL("kernel32.dll")
	procGetProcAddr = modKernel32.NewProc("GetProcAddress")
	themeMenu       *systray.MenuItem
	darkMenuEntry   *systray.MenuItem
	lightMenuEntry  *systray.MenuItem
)

func updateThemeMenuStatus(isDark bool) {
	if isDark {
		if darkMenuEntry != nil && !darkMenuEntry.Checked() {
			darkMenuEntry.Check()
			menuJsonData.IsDarkTheme = true
		}
		if lightMenuEntry != nil && lightMenuEntry.Checked() {
			lightMenuEntry.Uncheck()
			menuJsonData.IsLightTheme = false
		}
	} else {
		if darkMenuEntry != nil && darkMenuEntry.Checked() {
			darkMenuEntry.Uncheck()
			menuJsonData.IsDarkTheme = false
		}
		if lightMenuEntry != nil && !lightMenuEntry.Checked() {
			lightMenuEntry.Check()
			menuJsonData.IsLightTheme = true
		}
	}
	updateMenuJsonData()
}

func getProcByOrdinal(hModule syscall.Handle, ordinal uintptr) uintptr {
	ret, _, _ := procGetProcAddr.Call(uintptr(hModule), ordinal)
	return ret
}

func enableLightMode() {
	uxtheme, err := syscall.LoadLibrary("uxtheme.dll")
	if err != nil {
		return
	}
	pSetPreferredAppMode := getProcByOrdinal(uxtheme, 135)
	if pSetPreferredAppMode != 0 {
		syscall.SyscallN(pSetPreferredAppMode, 3)
	} else {
		pAllowDarkModeForApp := getProcByOrdinal(uxtheme, 132)
		if pAllowDarkModeForApp != 0 {
			syscall.SyscallN(pAllowDarkModeForApp, 0)
		}
	}
	pFlushMenuThemes := getProcByOrdinal(uxtheme, 136)
	if pFlushMenuThemes != 0 {
		syscall.SyscallN(pFlushMenuThemes)
	}
}

func enableDarkMode() {
	uxtheme, err := syscall.LoadLibrary("uxtheme.dll")
	if err != nil {
		return
	}
	pSetPreferredAppMode := getProcByOrdinal(uxtheme, 135)
	if pSetPreferredAppMode != 0 {
		syscall.SyscallN(pSetPreferredAppMode, 2)
	} else {
		pAllowDarkModeForApp := getProcByOrdinal(uxtheme, 132)
		if pAllowDarkModeForApp != 0 {
			syscall.SyscallN(pAllowDarkModeForApp, 1)
		}
	}
	pFlushMenuThemes := getProcByOrdinal(uxtheme, 136)
	if pFlushMenuThemes != 0 {
		syscall.SyscallN(pFlushMenuThemes)
	}
}

func buildThemeMenu(settingsMenu *systray.MenuItem) {
	themeMenu = settingsMenu.AddSubMenuItem("Theme", "Theme")
	darkMenuEntry = themeMenu.AddSubMenuItemCheckbox("Enable Dark", "Enable Dark", menuJsonData.IsDarkTheme)
	lightMenuEntry = themeMenu.AddSubMenuItemCheckbox("Enable Light", "Enable Light", menuJsonData.IsLightTheme)
	darkMenuEntry.Click(func() {
		enableDarkMode()
		updateThemeMenuStatus(true)
	})
	lightMenuEntry.Click(func() {
		enableLightMode()
		updateThemeMenuStatus(false)
	})
	if menuJsonData.IsDarkTheme {
		enableDarkMode()
	} else {
		enableLightMode()
	}
}
