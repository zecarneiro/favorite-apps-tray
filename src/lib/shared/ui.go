package shared

import (
	"errors"
	"fmt"
	"golangutils"
	"golangutils/entity"
	"main/src/entities"
	"strings"
)

func setUiArgs(funcName string, appId string, title string, message string, icon string) []string {
	args := ""
	if len(appId) > 0 {
		if SystemUtils.IsLinux() {
			args = "'%s'"
		} else if SystemUtils.IsWindows() {
			args = "-appId '%s'"
		}
		args = fmt.Sprintf(args, appId)
	}
	if len(title) > 0 {
		if SystemUtils.IsLinux() {
			args += " '%s'"
		} else if SystemUtils.IsWindows() {
			args += " -title '%s'"
		}
		args = fmt.Sprintf(args, title)
	}
	if len(message) > 0 {
		if SystemUtils.IsLinux() {
			args += " '%s'"
		} else if SystemUtils.IsWindows() {
			args += " -message '%s'"
		}
		args = fmt.Sprintf(args, message)
	}
	if len(icon) > 0 {
		if SystemUtils.IsLinux() {
			args += " '%s'"
		} else if SystemUtils.IsWindows() {
			args += " -icon '%s'"
		}
		args = fmt.Sprintf(args, icon)
	}
	return []string{"; " + funcName, args}
}

func getUiScriptCmd() entity.Command {
	cmd := entity.Command{
		Cmd:     ". " + golangutils.ResolvePath(VendorDir+"/powershell-utils/others/profile-shell/ui-powershell-functions.ps1"),
		Verbose: false,
	}
	if SystemUtils.IsLinux() {
		cmd.Cmd = "source '" + golangutils.ResolvePath(VendorDir+"/bash-utils/others/profile-shell/ui-bash-functions.sh") + "'"
		cmd.UseBash = true
	} else if SystemUtils.IsWindows() {
		cmd.UsePowerShell = true
	}
	return cmd
}

func Notify(title string, message string) {
	cmd := getUiScriptCmd()
	cmd.Args = setUiArgs("notify", GetAppNameFormated(), title, message, GetIcon())
	ConsoleUtils.ExecAsync(cmd, nil)
	if EnableLogs {
		LoggerUtils.Log(message)
	}
}

func OkNotify(message string) {
	cmd := getUiScriptCmd()
	cmd.Args = setUiArgs("oknotify", GetAppNameFormated(), "", message, GetIcon())
	ConsoleUtils.ExecAsync(cmd, nil)
	if EnableLogs {
		LoggerUtils.Ok(message)
	}
}

func InfoNotify(message string) {
	cmd := getUiScriptCmd()
	cmd.Args = setUiArgs("infonotify", GetAppNameFormated(), "", message, GetIcon())
	ConsoleUtils.ExecAsync(cmd, nil)
	if EnableLogs {
		LoggerUtils.Info(message)
	}
}

func WarnNotify(message string) {
	cmd := getUiScriptCmd()
	cmd.Args = setUiArgs("warnnotify", GetAppNameFormated(), "", message, GetIcon())
	ConsoleUtils.ExecAsync(cmd, nil)
	if EnableLogs {
		LoggerUtils.Warn(message)
	}
}

func ErrorNotify(message string) {
	cmd := getUiScriptCmd()
	cmd.Args = setUiArgs("errornotify", GetAppNameFormated(), "", message, GetIcon())
	ConsoleUtils.ExecAsync(cmd, nil)
	if EnableLogs {
		LoggerUtils.Error(message)
	}
}

func SelectFileDialog() (string, error) {
	cmd := getUiScriptCmd()
	cmd.Args = setUiArgs("selectfiledialog", "", "", "", "")
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
	cmd := entity.Command{
		Cmd:     "",
		Verbose: false,
		Args:    []string{},
	}
	if SystemUtils.IsLinux() {
		cmd.Cmd = "zenity"
		cmd.UseBash = true
		cmd.Args = append(
			cmd.Args,
			"--info",
			fmt.Sprintf("--text=\"%s\"", message),
			fmt.Sprintf("--title=\"%s\"", GetAppNameFormated()),
			fmt.Sprintf("--icon=\"%s\"", GetIcon()),
		)
	} else if SystemUtils.IsWindows() {
		cmd.Cmd = fmt.Sprintf("[System.Windows.Forms.MessageBox]::Show('%s', '%s', [System.Windows.Forms.MessageBoxButtons]::OK, '%s')", message, GetAppNameFormated(), GetIcon())
		cmd.UsePowerShell = true
	}
	if len(cmd.Cmd) > 0 {
		ConsoleUtils.ExecAsync(cmd, nil)
	}
}
