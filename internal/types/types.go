package types

// Discord webhook embed
type DiscordEmbedFooter struct {
	Text string `json:"text"`
}

// Discord webhook embed thumbnail
type DiscordEmbedThumbnail struct {
	Url string `json:"url"`
}

// Discord webhook embed
type DiscordEmbed struct {
	Title       string                `json:"title"`
	Description string                `json:"description"`
	Url         string                `json:"url"`
	Color       string                `json:"color"`
	Footer      DiscordEmbedFooter    `json:"footer"`
	TimeStamp   string                `json:"timestamp"`
	Thumbnail   DiscordEmbedThumbnail `json:"thumbnail"`
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

// App type
type App struct {
	Name          string
	Id            string
	Version       int
	DockerVersion string
}

// App status
type RuntipiAppStatus struct {
	Id      string `json:"id"`
	Version int    `json:"version"`
}

// App info
type RuntipiAppInfo struct {
	Name string `json:"name"`
}

// App update info
type RuntipiAppUpdateInfo struct {
	LatestVersion       int    `json:"latestVersion"`
	LatestDockerVersion string `json:"latestDockerVersion"`
}

// Runtipi App
type RuntipiApp struct {
	App        RuntipiAppStatus     `json:"app"`
	Info       RuntipiAppInfo       `json:"info"`
	UpdateInfo RuntipiAppUpdateInfo `json:"updateInfo"`
}

// Get apps response
type GetInstalledAppsResponse struct {
	Installed []RuntipiApp `json:"installed"`
}

// Server config
type ServerConfig struct {
	NotifyUrl          string `validate:"required" mapstructure:"notify-url"`
	RuntipiInternalUrl string `validate:"required" mapstructure:"runtipi-internal"`
	RuntipiUrl         string `validate:"required" mapstructure:"runtipi"`
	JwtSecret          string `validate:"required" mapstructure:"jwt-secret"`
	Appstore           string `validate:"required" mapstructure:"appstore"`
	DbPath             string `validate:"required" mapstructure:"db-path"`
	Refresh            int    `validate:"required" mapstructure:"refresh"`
	LogLevel           string `validate:"required,oneof=trace debug info warn error fatal panic" mapstructure:"log-level"`
	NoTls              bool   `validate:"required" mapstructure:"no-tls"`
}

// Check config
type CheckConfig struct {
	RuntipiUrl string `validate:"required" mapstructure:"runtipi"`
	JwtSecret  string `validate:"required" mapstructure:"jwt-secret"`
}
