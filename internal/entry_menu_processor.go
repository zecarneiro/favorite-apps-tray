package internal

import (
	"favoriteappstray/internal/shared"
	"fmt"
	"golangutils/pkg/common"
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
	message := fmt.Sprintf("Name: %s", shared.AppName)
	message += fmt.Sprintf("%sVersion: %s", common.Eol(), shared.AppVersion)
	message += fmt.Sprintf("%sRelease Date: %s", common.Eol(), shared.AppReleaseDate)
	message += fmt.Sprintf("%sLog file located: %s", common.Eol(), shared.GetLogFile())
	shared.InfoDialog(message)
}
