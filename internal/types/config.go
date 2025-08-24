package types

// API config
type APIConfig struct {
	RuntipiUrl string
	Secret     string
	Insecure   bool
}

// Alerts config
type AlertsConfig struct {
	NotificationUrl string
	RuntipiUrl      string
	Insecure        bool
	ServerName      string
}

// Server config
type ServerConfig struct {
	NotificationUrl string `validate:"required" mapstructure:"notification-url"`
	RuntipiUrl      string `validate:"required" mapstructure:"runtipi-url"`
	JwtSecret       string `validate:"required" mapstructure:"jwt-secret"`
	DatabasePath    string `mapstructure:"database-path"`
	Interval        int    `mapstructure:"interval"`
	LogLevel        string `validate:"oneof=trace debug info warn error fatal panic" mapstructure:"log-level"`
	Insecure        bool   `mapstructure:"insecure"`
	ServerName      string `mapstructure:"server-name"`
}

// Check config
type CheckConfig struct {
	RuntipiUrl string `validate:"required" mapstructure:"runtipi-url"`
	JwtSecret  string `validate:"required" mapstructure:"jwt-secret"`
	Insecure   bool   `mapstructure:"insecure"`
}
