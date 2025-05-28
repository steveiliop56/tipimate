package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"tipimate/internal/alerts"
	"tipimate/internal/api"
	"tipimate/internal/constants"
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
		handleError(err, "Failed to parse config")

		// Validate config
		err = validator.New().Struct(config)
		handleError(err, "Failed to validate config")

		// Configure logger
		log.Logger = log.Level(utils.GetLogLevel(config.LogLevel))
		log.Info().Str("version", constants.Version).Msg("Starting tipimate")

		// Validate URL
		sr := router.ServiceRouter{}
		_, err = sr.Locate(config.NotificationUrl)
		handleError(err, "Invalid notification URL")

		// Validate runtipi URL
		_, err = url.Parse(config.RuntipiUrl)
		handleError(err, "Invalid runtipi URL")

		// Initialize database
		db, err := database.InitDatabase(config.DatabasePath)
		handleError(err, "Failed to initialize database")

		// Migrate app urns
		log.Info().Msg("Migrating app URNs in the database")

		var dbApps []database.Schema
		db.Find(&dbApps)

		for _, dbApp := range dbApps {
			if !strings.Contains(dbApp.Urn, ":") {
				log.Info().Str("urn", dbApp.Urn).Msg("Migrating app URN to include appstore slug")

				// Migrate URN to include appstore slug
				db.Model(&dbApp).Updates(database.Schema{Urn: fmt.Sprintf("%s:migrated", dbApp.Urn)})
			}
		}

		// Create API
		apiConfig := types.APIConfig{
			RuntipiUrl: config.RuntipiUrl,
			Secret:     config.JwtSecret,
			Insecure:   config.Insecure,
		}

		api, err := api.NewAPI(apiConfig)
		handleError(err, "Failed to create API client")

		// Create alerts
		alertsConfig := types.AlertsConfig{
			NotificationUrl: config.NotificationUrl,
			RuntipiUrl:      config.RuntipiUrl,
			Insecure:        config.Insecure,
		}

		alerts := alerts.NewAlerts(alertsConfig)

		// Main loop
		for {
			log.Info().Msg("Checking for updates")

			// Get apps
			log.Info().Msg("Getting installed apps")
			apps, err := api.GetInstalledApps()
			handleError(err, "Failed to get installed apps")

			// Get appstores
			log.Info().Msg("Getting appstores")
			appstores, err := api.GetAppstores()
			handleError(err, "Failed to get appstores")

			// Get app ids
			installedApps := make(map[string]bool)
			for _, app := range apps.Installed {
				installedApps[app.Info.Urn] = true
			}

			// Get apps from database
			var dbApps []database.Schema
			db.Find(&dbApps)

			// Delete uninstalled apps
			for _, dbApp := range dbApps {
				if !installedApps[dbApp.Urn] {
					log.Warn().Str("urn", dbApp.Urn).Msg("Deleting app from the database")
					db.Unscoped().Delete(&dbApp)
				}
			}

			// Get apps with updates
			log.Info().Msg("Comparing versions")
			appsWithUpdates := []types.App{}

			for _, app := range apps.Installed {
				// If app is up to date, ignore it
				if app.App.Version == app.Metadata.LatestVersion {
					continue
				}

				log.Debug().Interface("app", app).Msg("App has an update")

				// Get app from database
				var dbApp database.Schema
				dbRes := db.First(&dbApp, "urn = ?", app.Info.Urn)

				// Check if app is not in database
				if dbRes.RowsAffected == 0 {
					log.Debug().Str("urn", app.Info.Urn).Msg("App not found in database, creating new entry")

					// Create app in database
					db.Create(&database.Schema{Urn: app.Info.Urn, Version: app.App.Version, LatestVersion: app.Metadata.LatestVersion})

					// Add app to updates
					appsWithUpdates = append(appsWithUpdates, types.App{
						Urn:           app.Info.Urn,
						Name:          app.Info.Name,
						Version:       app.App.Version,
						DockerVersion: app.Metadata.LatestDockerVersion,
					})
				} else {
					log.Debug().Str("urn", app.Info.Urn).Msg("App found in database, checking versions")

					// Modify db if version is different
					if dbApp.Version != app.App.Version || dbApp.LatestVersion != app.Metadata.LatestVersion {
						log.Debug().Str("urn", app.Info.Urn).Msg("Updating app in database")

						db.Model(&dbApp).Updates(database.Schema{LatestVersion: app.Metadata.LatestVersion, Version: app.App.Version})

						// Add app to updates
						appsWithUpdates = append(appsWithUpdates, types.App{
							Urn:           app.Info.Urn,
							Name:          app.Info.Name,
							Version:       app.App.Version,
							DockerVersion: app.Metadata.LatestDockerVersion,
						})
					}
				}
			}

			// Send notifications
			log.Info().Msg("Sending notifications")

			for _, appWithUpdate := range appsWithUpdates {
				// Log
				log.Logger.Info().Str("urn", appWithUpdate.Urn).Str("tipiVersion", strconv.Itoa(appWithUpdate.Version)).Str("dockerVersion", appWithUpdate.DockerVersion).Msg("App has an update")

				// Send alert
				alertErr := alerts.SendAlert(&appWithUpdate, appstores.Appstores)

				// Handle error
				handleError(alertErr, "Failed to send alert")
			}

			// Sleep
			time.Sleep(time.Duration(config.Interval) * time.Minute)
		}

	},
}

// Handle error
func handleError(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

// Add command
func init() {
	// Viper
	viper.SetEnvPrefix("tipimate")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Flags
	serverCmd.Flags().String("notification-url", "", "Notification URL (shoutrrr format)")
	serverCmd.Flags().String("runtipi-url", "", "Runtipi server URL")
	serverCmd.Flags().String("jwt-secret", "", "JWT secret")
	serverCmd.Flags().String("database-path", "tipimate.db", "Database path")
	serverCmd.Flags().Int("interval", 30, "Refresh interval in minutes")
	serverCmd.Flags().String("log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic)")
	serverCmd.Flags().Bool("insecure", true, "Disable TLS (https) for services like Gotify, Ntfy etc.")

	// Bind flags to viper
	viper.BindPFlags(serverCmd.Flags())

	// Add command
	rootCmd.AddCommand(serverCmd)
}
