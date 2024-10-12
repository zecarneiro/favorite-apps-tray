package shared

import "github.com/zecarneiro/golangutils"

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
		golangutils.Notify(APP_NAME, APP_NAME, message, icon)
	} else if typeNotify == 1 {
		golangutils.OkNotify(APP_NAME, message, icon)
	} else if typeNotify == 2 {
		golangutils.InfoNotify(APP_NAME, message, icon)
	} else if typeNotify == 3 {
		golangutils.WarnNotify(APP_NAME, message, icon)
	} else if typeNotify == 4 {
		golangutils.ErrorNotify(APP_NAME, message, icon)
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
