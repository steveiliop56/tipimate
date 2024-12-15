package types

// Discord webhook embed
type DiscordEmbedFooter struct {
	Text string `json:"text"`
}

type DiscordEmbedThumbnail struct {
	Url string `json:"url"`
}

type DiscordEmbed struct {
	Title string `json:"title"`
	Description	string `json:"description"`
	Url	string `json:"url"`
	Color string `json:"color"`
	Footer DiscordEmbedFooter `json:"footer"`
	TimeStamp string `json:"timestamp"`
	Thumbnail DiscordEmbedThumbnail `json:"thumbnail"`
}

// Discord webhook message
type DiscordMessage struct {
	Embeds []DiscordEmbed `json:"embeds"`
	AvatarUrl string `json:"avatar_url"`
	Username string `json:"username"`
}

// Discord webhook struct
type DiscordWebhook struct {
	Json bool `url:"json"`
}

// Ntfy webhook struct
type NtfyWebhook struct {
	Click string `url:"click"`
	Icon string `url:"icon"`
	Title string `url:"title"`
	Scheme string `url:"scheme"`
}

// App type
type SimpleApp struct {
	Name string
	Id string
	Version int
	DockerVersion string
}

// App update alert
type AppUpdateAlert struct {
	Name string
	Id string
	Version int
	DockerVersion string
	ServerUrl string
	DiscordUrl string
	AppStore string
}

// App status
type AppStatus struct {
	Id string `json:"id"`
	Version int `json:"version"`
}

// App info
type AppInfo struct {
	Name string `json:"name"`
}

// App update info
type AppUpdateInfo struct {
	LatestVersion int `json:"latestVersion"`
	LatestDockerVersion string `json:"latestDockerVersion"`
}

// App
type App struct {
	App AppStatus `json:"app"`
	Info AppInfo `json:"info"`
	UpdateInfo AppUpdateInfo `json:"updateInfo"`
}

// Get apps response
type GetInstalledAppsResponse struct {
	Installed []App `json:"installed"`
}

// Config
type ServerConfig struct {
	NotifyUrl string `validate:"required" message:"Notify URL is required" mapstructure:"notify-url"`
	RuntipiInternalUrl string `validate:"required" message:"Runtipi internal URL is required" mapstructure:"runtipi-internal"`
	RuntipiUrl string `validate:"required" message:"Runtipi URL is required" mapstructure:"runtipi"`
	JwtSecret string `validate:"required" message:"JWT secret is required" mapstructure:"jwt-secret"`
	Appstore string `validate:"required" message:"Appstore URL is required" mapstructure:"appstore"`
	DbPath string `validate:"required" message:"Database path is required" mapstructure:"db-path"`
	Refresh int `validate:"required" message:"Refresh interval is required" mapstructure:"refresh"`
}

type CheckConfig struct {
	RuntipiUrl string `validate:"required" message:"Runtipi URL is required" mapstructure:"runtipi"`
	JwtSecret string `validate:"required" message:"JWT secret is required" mapstructure:"jwt-secret"`
}