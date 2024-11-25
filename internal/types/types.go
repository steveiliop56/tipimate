package types

// Discord webhook embed
type EmbedFooter struct {
	Text string `json:"text"`
}

type EmbedThumbnail struct {
	Url string `json:"url"`
}

type Embed struct {
	Title string `json:"title"`
	Description	string `json:"description"`
	Url	string `json:"url"`
	Color string `json:"color"`
	Footer EmbedFooter `json:"footer"`
	TimeStamp string `json:"timestamp"`
	Thumbnail EmbedThumbnail `json:"thumbnail"`
}

// Discord webhook message
type Message struct {
	Embeds []Embed `json:"embeds"`
	AvatarUrl string `json:"avatar_url"`
	Username string `json:"username"`
}

// Discord webhook struct
type Webhook struct {
	Json bool `url:"json"`
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