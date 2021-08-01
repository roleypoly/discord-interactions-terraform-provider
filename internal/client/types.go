package client

type InteractionCommand struct {
	ID                string `json:"id,omitempty"`
	ApplicationID     string `json:"application_id,omitempty"`
	GuildID           string `json:"guild_id,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	DefaultPermission bool   `json:"default_permission,omitempty"`

	Options []InteractionCommandOption `json:"options,omitempty"`
}

type InteractionCommandOption struct {
	Type        int    `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`

	Choices []InteractionCommandOptionChoice `json:"choices,omitempty"`
	Options []InteractionCommandOption       `json:"options,omitempty"`
}

type InteractionCommandOptionChoice struct {
	Name string `json:"name,omitempty"`

	// Value can be a string, int, or float
	Value interface{} `json:"value,omitempty"`
}
