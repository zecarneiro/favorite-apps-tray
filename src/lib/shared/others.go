package shared

import (
	"errors"
	"golangutils"
	"golangutils/entity"
	"main/src/entities"
	"strings"
)

func getUiScriptCmd() entity.Command {
	cmd := entity.Command{
		Cmd:     ". '" + golangutils.JoinPath(ScriptsDir, "powershell-functions.ps1") + "'",
		Verbose: EnableLogs,
	}
	if SystemUtils.IsLinux() {
		cmd.Cmd = "source '" + golangutils.JoinPath(ScriptsDir, "bash-functions.sh") + "'"
		cmd.Args = []string{SystemUtils.Info().PlatformName}
		cmd.UseBash = true
	} else if SystemUtils.IsWindows() {
		cmd.UsePowerShell = true
	}
	return cmd
}

func Notify(title string, message string, noLog bool) {
	cmd := getUiScriptCmd()
	cmd.Args = []string{"; notify"}
	appIdKey := "FAT_NOTIFY_APP_ID"
	titleKey := "FAT_NOTIFY_TITLE"
	messageKey := "FAT_NOTIFY_MESSAGE"
	iconKey := "FAT_NOTIFY_ICON"

	golangutils.SetEnv(appIdKey, GetAppNameFormated())
	golangutils.SetEnv(titleKey, title)
	golangutils.SetEnv(messageKey, message)
	golangutils.SetEnv(iconKey, GetIcon())
	ConsoleUtils.ExecAsync(cmd, ProcessConsoleResult)
	if EnableLogs && !noLog {
		LoggerUtils.Log(message)
	}
	golangutils.UnsetEnv(appIdKey)
	golangutils.UnsetEnv(titleKey)
	golangutils.UnsetEnv(messageKey)
	golangutils.UnsetEnv(iconKey)
}

func OkNotify(message string) {
	Notify("Success", message, true)
	if EnableLogs {
		LoggerUtils.Ok(message)
	}
}

func InfoNotify(message string) {
	Notify("Information", message, true)
	if EnableLogs {
		LoggerUtils.Info(message)
	}
}

func WarnNotify(message string) {
	Notify("Warning", message, true)
	if EnableLogs {
		LoggerUtils.Warn(message)
	}
}

func ErrorNotify(message string) {
	Notify("Error", message, true)
	if EnableLogs {
		LoggerUtils.Error(message)
	}
}

func SelectFileDialog() (string, error) {
	cmd := getUiScriptCmd()
	cmd.Args = []string{"; selectfiledialog"}
	res := ConsoleUtils.Exec(cmd)
	if res.HasError() {
		return "", res.Error
	}
	if SystemUtils.IsWindows() {
		selected, err := golangutils.StringToObject[entities.SelectedFileDialog](res.Data)
		if err != nil {
			return "", err
		}
		return selected.Selected, err
	} else if SystemUtils.IsLinux() {
		resArr := strings.Split(res.Data, SystemUtils.Info().Eol)
		if len(resArr) == 0 {
			return "", errors.New("No file selected")
		}

		if len(resArr) > 1 {
			res.Data = resArr[1]
		} else {
			res.Data = resArr[0]
		}
		selected, err := golangutils.StringToObject[entities.SelectedFileDialog](res.Data)
		return selected.Selected, err
	}
	return "", errors.New(golangutils.GetUnknowOSMsg())
}

func ShowMessageDialog(message string) {
	cmd := getUiScriptCmd()
	cmd.Args = []string{"; messagedialog"}
	titleKey := "FAT_DIALOG_TITLE"
	messageKey := "FAT_DIALOG_MESSAGE"
	iconKey := "FAT_DIALOG_ICON"

	golangutils.SetEnv(titleKey, GetAppNameFormated())
	golangutils.SetEnv(messageKey, message)
	golangutils.SetEnv(iconKey, GetIcon())
	ConsoleUtils.ExecAsync(cmd, ProcessConsoleResult)
	golangutils.UnsetEnv(titleKey)
	golangutils.UnsetEnv(messageKey)
	golangutils.UnsetEnv(iconKey)
}

func ExtractIcon(file string, dest string) bool {
	cmd := getUiScriptCmd()
	cmd.Args = []string{"; iconextractor"}
	fileKey := "FAT_ICON_EXTRACTOR_FILE"
	destKey := "FAT_ICON_EXTRACTOR_DEST"

	golangutils.SetEnv(fileKey, file)
	golangutils.SetEnv(destKey, dest)
	ProcessConsoleResult(ConsoleUtils.Exec(cmd))
	golangutils.UnsetEnv(fileKey)
	golangutils.UnsetEnv(destKey)
	if golangutils.FileExist(dest) {
		return true
	}
	return false
}

func AppsInfo(typeApp string) {
	cmd := getUiScriptCmd()
	cmd.Args = []string{"; appinfo"}
	appNameKey := "FAT_INFO_APP_NAME"
	typeAppKey := "FAT_INFO_TYPE"

	golangutils.SetEnv(appNameKey, ApplicationName)
	golangutils.SetEnv(typeAppKey, typeApp)
	ProcessConsoleResult(ConsoleUtils.Exec(cmd))
	golangutils.UnsetEnv(appNameKey)
	golangutils.UnsetEnv(typeAppKey)
}

func ShowProcessingMsg(isDone bool) {
	if isDone {
		OkNotify("Processing, done.")
	} else {
		InfoNotify("Processing...")
	}
}
