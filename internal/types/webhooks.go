package types

// Discord webhook embed
type DiscordEmbedFooter struct {
	Text string `json:"text"`
}

// Discord webhook embed
type DiscordEmbed struct {
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Url         string             `json:"url"`
	Color       string             `json:"color"`
	Footer      DiscordEmbedFooter `json:"footer"`
	Timestamp   string             `json:"timestamp"`
}

// Discord webhook message
type DiscordMessage struct {
	Embeds    []DiscordEmbed `json:"embeds"`
	AvatarUrl string         `json:"avatar_url"`
	Username  string         `json:"username"`
}

// Discord webhook struct
type DiscordWebhook struct {
	Json bool `url:"json"`
}

// Ntfy webhook struct
type NtfyWebhook struct {
	Click  string `url:"click"`
	Icon   string `url:"icon"`
	Title  string `url:"title"`
	Scheme string `url:"scheme"`
}

// Gotify webhook struct
type GotifyWebhook struct {
	DisableTls bool   `url:"disableTls"`
	Title      string `url:"title"`
}
