package cmd

import (
	"net/url"
	"os"
	"strconv"
	"time"
	"tipimate/internal/alerts"
	"tipimate/internal/api"
	"tipimate/internal/assets"
	"tipimate/internal/database"
	"tipimate/internal/types"
	"tipimate/internal/utils"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the tipimate server",
	Long:  "Use the server command to automatically check for updates on your runtipi server and send you notifications when updates are available",
	Run: func(cmd *cobra.Command, args []string) {
		// Logger
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger().Level(zerolog.FatalLevel)

		// Get config
		var config types.ServerConfig
		err := viper.Unmarshal(&config)
		handleErrorLogger(err, "Failed to parse config")

		if config.RuntipiInternalUrl == "" {
			config.RuntipiInternalUrl = config.RuntipiUrl
		}

		// Validate config
		err = validator.New().Struct(config)
		handleErrorLogger(err, "Failed to validate config")

		// Configure logger
		log.Level(utils.GetLogLevel(config.LogLevel))
		log.Info().Str("version", assets.Version).Msg("Starting tipimate server")

		// Validate URL
		sr := router.ServiceRouter{}
		_, err = sr.Locate(config.NotifyUrl)
		handleErrorLogger(err, "Invalid notification URL")

		// Validate runtipi URL
		_, err = url.Parse(config.RuntipiUrl)
		handleErrorLogger(err, "Invalid runtipi URL")

		// Initialize db
		db, err := database.InitDb(config.DbPath)
		handleErrorLogger(err, "Failed to initialize database")

		// Create JWT
		token, err := api.CreateJWT(config.JwtSecret)
		handleErrorLogger(err, "Failed to create JWT token")

		// Main loops
		for {
			log.Info().Msg("Checking for updates")

			// Get installed apps
			log.Info().Msg("Getting installed apps")
			installedApps, err := api.GetInstalledApps(token, config.RuntipiInternalUrl)
			handleErrorLogger(err, "Failed to get installed apps")

			// Get installed app ids
			installedAppIds := make(map[string]bool)
			for _, app := range installedApps.Installed {
				installedAppIds[app.App.Id] = true
			}

			// Create db apps
			var dbApps []database.Schema
			db.Find(&dbApps)

			// Delete uninstalled apps
			for _, dbApp := range dbApps {
				if !installedAppIds[dbApp.Id] {
					log.Warn().Str("appId", dbApp.Id).Msg("Deleting app from the database")
					db.Unscoped().Delete(&dbApp)
				}
			}

			// Get apps with updates
			log.Info().Msg("Comparing versions")
			appsWithUpdates := []types.App{}

			for _, app := range installedApps.Installed {
				// Get app from database
				var dbApp database.Schema
				dbRes := db.First(&dbApp, "id = ?", app.App.Id)

				// Check if app is not in database
				if dbRes.RowsAffected == 0 {
					// Create app in database
					db.Create(&database.Schema{Id: app.App.Id, Version: app.App.Version, LatestVersion: app.UpdateInfo.LatestVersion})

					// Check for updates
					if app.App.Version < app.UpdateInfo.LatestVersion {
						appsWithUpdates = append(appsWithUpdates, types.App{
							Id:            app.App.Id,
							Name:          app.Info.Name,
							Version:       app.App.Version,
							DockerVersion: app.UpdateInfo.LatestDockerVersion,
						})
					}
				} else {
					// Check if app has changed
					if dbApp.LatestVersion != app.UpdateInfo.LatestVersion || dbApp.Version != app.App.Version {
						// Check for updates
						if app.App.Version < app.UpdateInfo.LatestVersion {
							appsWithUpdates = append(appsWithUpdates, types.App{
								Id:            app.App.Id,
								Name:          app.Info.Name,
								Version:       app.App.Version,
								DockerVersion: app.UpdateInfo.LatestDockerVersion,
							})
						}

						// Modify db
						db.Model(&dbApp).Updates(database.Schema{LatestVersion: app.UpdateInfo.LatestVersion, Version: app.App.Version})
					}
				}
			}

			// Send notifications
			log.Info().Msg("Sending notifications")
			for _, appWithUpdate := range appsWithUpdates {
				// Log
				log.Logger.Info().Str("appId", appWithUpdate.Id).Str("tipiVersion", strconv.Itoa(appWithUpdate.Version)).Str("dockerVersion", appWithUpdate.DockerVersion).Msg("App has an update")

				// Send alert
				alertErr := alerts.SendAlert(&appWithUpdate, config.NotifyUrl, config.RuntipiUrl, config.Appstore, config.NoTls)

				// Handle error
				handleErrorLogger(alertErr, "Failed to send alert")
			}

			// Sleep
			time.Sleep(time.Duration(config.Refresh) * time.Minute)
		}

	},
}

// Handle error with logger
func handleErrorLogger(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

// Add command
func init() {
	// Viper
	viper.AutomaticEnv()

	// Flags
	serverCmd.Flags().String("notify-url", "", "Notification URL (shoutrrr format)")
	serverCmd.Flags().String("runtipi", "", "Runtipi server URL")
	serverCmd.Flags().String("runtipi-internal", "", "Runtipi internal URL (used when running in the same server as runtipi)")
	serverCmd.Flags().String("jwt-secret", "", "JWT secret")
	serverCmd.Flags().String("appstore", "https://github.com/runtipi/runtipi-appstore", "Appstore URL for images")
	serverCmd.Flags().String("db-path", "tipimate.db", "Database path")
	serverCmd.Flags().String("log-level", "info", "Log level")
	serverCmd.Flags().Int("refresh", 30, "Refresh interval")
	serverCmd.Flags().Bool("no-tls", true, "Disable TLS (https) for services like Gotify, Ntfy etc.")

	// Bind flags
	viper.BindEnv("notify-url", "NOTIFY_URL")
	viper.BindEnv("runtipi", "RUNTIPI_URL")
	viper.BindEnv("runtipi-internal", "RUNTIPI_INTERNAL")
	viper.BindEnv("jwt-secret", "JWT_SECRET")
	viper.BindEnv("appstore", "APPSTORE")
	viper.BindEnv("db-path", "DB_PATH")
	viper.BindEnv("log-level", "LOG_LEVEL")
	viper.BindEnv("no-tls", "NO_TLS")

	// Bind flags to viper
	viper.BindPFlags(serverCmd.Flags())

	// Add command
	rootCmd.AddCommand(serverCmd)
}
