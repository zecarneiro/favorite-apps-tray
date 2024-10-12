package shared

import (
	"errors"
	"main/src/entities"
	"main/src/lib/golangutils"
	"strings"
)

func getIcon() string {
	icon := ExecutableDir + "/assets/image/logo/icon"
	if golangutils.IsWindows() {
		icon = icon + ".ico"
	} else if golangutils.IsLinux() {
		icon = icon + ".png"
	}
	return golangutils.ResolvePath(icon)
}

func showNotification(message string, typeNotify int) {
	icon := getIcon()
	if typeNotify == 0 {
		golangutils.Notify(ApplicationName, ApplicationName, message, icon)
	} else if typeNotify == 1 {
		golangutils.OkNotify(ApplicationName, message, icon)
	} else if typeNotify == 2 {
		golangutils.InfoNotify(ApplicationName, message, icon)
	} else if typeNotify == 3 {
		golangutils.WarnNotify(ApplicationName, message, icon)
	} else if typeNotify == 4 {
		golangutils.ErrorNotify(ApplicationName, message, icon)
	}
}

func Notify(message string) {
	showNotification(message, 0)
}

func OkNotify(message string) {
	showNotification(message, 1)
}

func InfoNotify(message string) {
	showNotification(message, 2)
}

func WarnNotify(message string) {
	showNotification(message, 3)
}

func ErrorNotify(message string) {
	showNotification(message, 4)
}

func SelectFileDialog() (string, error) {
	cmd := golangutils.CommandInfo{
		Cmd:     ". " + golangutils.ResolvePath(VendorDir+"/powershell-utils/others/profile-shell/ui-powershell-functions.ps1"),
		Args:    []string{"; selectfiledialog"},
		Verbose: false,
	}
	if golangutils.IsLinux() {
		cmd.Cmd = "source '" + golangutils.ResolvePath(VendorDir+"/bash-utils/others/profile-shell/ui-bash-functions.sh") + "'"
		cmd.UseBash = true
	} else if golangutils.IsWindows() {
		cmd.UsePowerShell = true
	}
	res := golangutils.Exec(cmd)
	if res.HasError() {
		return "", res.Error
	}
	if golangutils.IsWindows() {
		selected, err := golangutils.StringToObject[entities.SelectedFileDialog](res.Data)
		if err!= nil {
            return "", err
        }
		return selected.Selected, err
	} else if golangutils.IsLinux() {
		resArr := strings.Split(res.Data, golangutils.SysInfo().Eol)
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
	return "", errors.New(golangutils.UNKNOW_OS_MSG)
}
