package internal

import (
	"favoriteappstray/internal/shared"
	"fmt"
	"golangutils/pkg/file"
	"golangutils/pkg/ui"
)

func updateMenuProcessor() {
	shared.ShowProcessingMsg(false)
	refresh(true)
	shared.ShowProcessingMsg(true)
}

func selectJsonFileProcessor() {
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
}

func enableLogsProcessor() {
	message := ""
	if enableLogsMenuItem.Checked() {
		enableLogsMenuItem.Uncheck()
		message = "disabled"
		shared.EnableLogs = false
	} else {
		enableLogsMenuItem.Check()
		message = "enabled"
		shared.EnableLogs = true
	}
	menuJsonData.EnableLogs = shared.EnableLogs
	updateMenuJsonData()
	shared.InfoNofity(fmt.Sprintf("All Logs was %s by user.", message))
}

func aboutProcessor() {
	message := "Name: " + shared.AppName
	message += "\nVersion: " + shared.AppVersion
	message += "\nRelease Date: " + shared.AppReleaseDate
	message += "\nLog file located: " + shared.GetLogFile()
	shared.InfoDialog(message)
}
