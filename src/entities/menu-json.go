package entities

type MenuJson struct {
	NoMenu     []MenuItemJson            `json:"noMenu,omitempty"`
	Others     map[string][]MenuItemJson `json:"others,omitempty"`
	EnableLogs bool                      `json:"enableLogs,omitempty"`
}
