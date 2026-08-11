package entities

type AppsInfo struct {
	DisplayName string `json:"displayName,omitempty"`
	Shortcut    string `json:"shortcut,omitempty"`
	Command     string `json:"command,omitempty"`
	Icon        string `json:"icon,omitempty"`
}
