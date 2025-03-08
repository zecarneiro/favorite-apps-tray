package shared

import (
	"golangutils"
	"golangutils/entity"
	"main/src/entities"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ExecutableDir, _       = golangutils.GetExecutableDir()
	ScriptsDir             = golangutils.ResolvePath(ExecutableDir + "/scripts")
	ApplicationName        string
	ApplicationVersion     string
	ApplicationReleaseDate string
	LoggerUtils            = golangutils.NewLoggerUtils()
	SystemUtils            = golangutils.NewSystemUtils(&LoggerUtils)
	ConsoleUtils           = golangutils.NewConsoleUtils(&LoggerUtils)
	EnableLogs             = true
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

func setScriptsPermission() {
	if SystemUtils.IsWindows() {
		ConsoleUtils.SetFullPermissionWindows(ExecutableDir)
	} else if SystemUtils.IsLinux() {
		ConsoleUtils.SetFullPermissionLinux(ExecutableDir)
	}
}

func GetAppNameFormated() string {
	nameFormated := golangutils.StringReplaceAll(ApplicationName, map[string]string{"-": " "})
	caser := cases.Title(language.English)
	return caser.String(strings.ToLower(nameFormated))
}

func GetConfigIcon(appName string) string {
	iconDir := golangutils.ResolvePath(GetConfigurationDir() + "/icon")
	golangutils.CreateDirectory(iconDir, true)
	return golangutils.ResolvePath(iconDir + "/" + appName)
}

func GetLogFile() string {
	return golangutils.JoinPath(GetConfigurationDir(), "log")
}

func GetConfigurationDir() string {
	configDir := golangutils.ResolvePath(SystemUtils.Info().HomeDir + "/.config/favorite-apps-tray")
	golangutils.CreateDirectory(configDir, true)
	return configDir
}

func GetJsonFile() string {
	return golangutils.ResolvePath(GetConfigurationDir() + "/tray_menu_entries.json")
}

func GetIcon() string {
	icon := ExecutableDir + "/assets/image/logo/icon"
	if SystemUtils.IsWindows() {
		icon = icon + ".ico"
	} else if SystemUtils.IsLinux() {
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
	logFile := GetLogFile()
	golangutils.DeleteFile(logFile)
	LoggerUtils.SetLogFile(logFile)
	golangutils.ReadFileLineByLine(golangutils.ResolvePath(ExecutableDir+"/APP_INFO.conf"), loadAppInformations)
	setScriptsPermission()
}

func ProcessConsoleResult(result entity.Response[string]) {
	if len(result.Data) > 0 {
		LoggerUtils.Debug(result.Data)
	}
	if result.HasError() {
		LoggerUtils.Error(result.Error.Error())
	}
}
