package entities

type MenuItemJson struct {
	Name          string `json:"name,omitempty"`
	Type          string `json:"type,omitempty"`
	IsInitialName bool   `json:"isInitialName,omitempty"`
	Command       string `json:"command,omitempty"`
}
