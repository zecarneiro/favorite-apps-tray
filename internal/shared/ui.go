package shared

import (
	"golangutils/pkg/logger"
	"golangutils/pkg/ui"
)

var isStartProcessingMsg = false

func InfoNofity(message string) error {
	if EnableLogs {
		logger.Info(message)
	}
	return ui.InfoNofity(message, GetIcon())
}

func WarnNofity(message string) error {
	if EnableLogs {
		logger.Warn(message)
	}
	return ui.WarnNofity(message, GetIcon())
}

func ErrorNofity(message string) error {
	if EnableLogs {
		logger.ErrorStr(message)
	}
	return ui.ErrorNofity(message, GetIcon())
}

func OkNofity(message string) error {
	if EnableLogs {
		logger.Ok(message)
	}
	return ui.OkNofity(message, GetIcon())
}

func ShowProcessingMsg(isDone bool) {
	if isDone {
		if isStartProcessingMsg {
			isStartProcessingMsg = false
			OkNofity("Processing, done.")
		}
	} else {
		if !isStartProcessingMsg {
			isStartProcessingMsg = true
			InfoNofity("Processing...")
		}
	}
}

func InfoDialog(message string) error {
	return ui.InfoDialog(AppName, message)
}

func WarnDialog(message string) error {
	return ui.WarnDialog(AppName, message)
}

func ErrorDialog(message string) error {
	return ui.ErrorDialog(AppName, message)
}
