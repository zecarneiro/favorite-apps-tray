package enums

type TypeApps string

var (
	WINDOWS_APPS TypeApps = "windows-apps"
	SHORTCUTS    TypeApps = "shortcuts"
	COMMAND      TypeApps = "command"
)

func (t TypeApps) String() string {
	return string(t)
}
