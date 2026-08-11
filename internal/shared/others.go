package shared

import (
	internalEnums "favoriteappstray/internal/enums"
	"fmt"
	"golangutils/pkg/env"
	"golangutils/pkg/exe"
	"golangutils/pkg/file"
	"golangutils/pkg/logger"
	"golangutils/pkg/logic"
	"golangutils/pkg/models"
	"golangutils/pkg/platform"
)

func ExtractIcon(fileSrc string, dest string) bool {
	if platform.IsLinux() && file.IsFile(fileSrc) {
		cmd := models.Command{
			Cmd:      fmt.Sprintf(`inkscape "%s" -o "%s" --export-overwrite -w 32 -h 32`, fileSrc, dest),
			UseShell: true,
			Verbose:  EnableLogs,
		}
		err := exe.ExecRealTime(cmd)
		if err != nil {
			logger.Error(err)
		}
	} else if platform.IsWindows() {
		fileKey := "FAT_ICON_EXTRACTOR_FILE"
		destKey := "FAT_ICON_EXTRACTOR_DEST"
		env.SetWithSingleValue(fileKey, fileSrc)
		env.SetWithSingleValue(destKey, dest)
		if err := exe.ExecRealTime(models.Command{Cmd: file.JoinPath(GetScriptsDir(), "icon-extractor-for-windows.ps1"), UseShell: true}); err != nil {
			logger.Error(err)
		}
		env.Unset(fileKey)
		env.Unset(destKey)
	}
	return logic.Ternary(file.IsFile(dest), true, false)
}

func AppsInfo(typeApp internalEnums.TypeApps) {
	typeAppKey := "FAT_INFO_TYPE"
	jsonFileWriter := func(data string) {
		jsonFileConfig := models.FileWriterConfig{
			File:        file.JoinPath(GetConfigurationDir(), fmt.Sprintf(`apps-info-%s.json`, typeApp.String())),
			Data:        data,
			IsAppend:    false,
			IsCreateDir: true,
		}
		file.WriteFile(jsonFileConfig)
	}
	if platform.IsLinux() {
		data, err := exe.Exec(models.Command{Cmd: fmt.Sprintf(`gjs %s`, file.JoinPath(GetScriptsDir(), "app-info-for-gnome.js")), UseShell: true})
		if err != nil {
			ErrorNofity(err.Error())
		} else {
			jsonFileWriter(data)
		}
	} else if platform.IsWindows() {
		env.SetWithSingleValue(typeAppKey, typeApp.String())
		data, err := exe.Exec(models.Command{Cmd: file.JoinPath(GetScriptsDir(), "app-info-for-windows.ps1"), UseShell: true})
		env.Unset(typeAppKey)
		if err != nil {
			ErrorNofity(err.Error())
		} else {
			jsonFileWriter(data)
		}
	}
}
