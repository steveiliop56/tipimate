package types

// App status
type RuntipiAppStatus struct {
	Version int `json:"version"`
}

// App info
type RuntipiAppInfo struct {
	Name string `json:"name"`
	Urn  string `json:"urn"`
}

// App update info
type RuntipiAppMetadata struct {
	LatestVersion       int    `json:"latestVersion"`
	LatestDockerVersion string `json:"latestDockerVersion"`
}

// Runtipi App
type RuntipiApp struct {
	App      RuntipiAppStatus   `json:"app"`
	Info     RuntipiAppInfo     `json:"info"`
	Metadata RuntipiAppMetadata `json:"metadata"`
}

// Get apps response
type GetInstalledAppsResponse struct {
	Installed []RuntipiApp `json:"installed"`
}

// Appstore
type RuntipiAppstore struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	Enabled bool   `json:"enabled"`
}

// Get appstores response
type GetAppstoresResponse struct {
	Appstores []RuntipiAppstore `json:"appStores"`
}
