package entities

type MenuItemJson struct {
	Name               string `json:"name,omitempty"`
	Type               string `json:"type,omitempty"`
	Command            string `json:"command,omitempty"`
	Regex              string `json:"regex,omitempty"`
	RegexOnDisplayName bool   `json:"regexOnDisplayName,omitempty"`
}
