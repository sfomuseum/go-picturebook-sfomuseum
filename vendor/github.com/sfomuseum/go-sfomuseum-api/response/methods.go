package response

type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Example     string `json:"example"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Method struct {
	Name                string      `json:"name"`
	RequestMethod       string      `json:"request_method"`
	Description         string      `json:"description"`
	RequiresPermissions int         `json:"requires_perms"`
	Parameters          []Parameter `json:"parameters"`
	Errors              []Error     `json:"errors"`
	Notes               []string    `json:"notes"`
	Extras              bool        `json:"extras"`
	Paginated           bool        `json:"paginated"`
}
