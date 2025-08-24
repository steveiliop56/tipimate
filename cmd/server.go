package cmd

import (
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

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the tipimate server",
	Long:  "Use the server command to automatically check for updates on your runtipi server and send you notifications when updates are available",
	Run: func(cmd *cobra.Command, args []string) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger().Level(zerolog.FatalLevel)

		var config types.ServerConfig
		err := viper.Unmarshal(&config)
		handleError(err, "Failed to parse config")

		err = validator.New().Struct(config)
		handleError(err, "Failed to validate config")

		log.Logger = log.Level(utils.GetLogLevel(config.LogLevel))
		log.Info().Str("version", constants.Version).Msg("Starting tipimate")

		log.Debug().Interface("config", config).Msg("Dumping configuration")

		sr := router.ServiceRouter{}
		_, err = sr.Locate(config.NotificationUrl)
		handleError(err, "Invalid notification URL")

		_, err = url.Parse(config.RuntipiUrl)
		handleError(err, "Invalid runtipi URL")

		db, err := database.InitDatabase(config.DatabasePath)
		handleError(err, "Failed to initialize database")

		apiConfig := types.APIConfig{
			RuntipiUrl: config.RuntipiUrl,
			Secret:     config.JwtSecret,
			Insecure:   config.Insecure,
		}

		api, err := api.NewAPI(apiConfig)
		handleError(err, "Failed to create API client")

		alertsConfig := types.AlertsConfig{
			NotificationUrl: config.NotificationUrl,
			RuntipiUrl:      config.RuntipiUrl,
			Insecure:        config.Insecure,
			ServerName:      config.ServerName,
		}

		alerts := alerts.NewAlerts(alertsConfig)

		for range time.Tick(time.Duration(config.Interval) * time.Minute) {
			log.Info().Msg("Checking for updates")

			log.Info().Msg("Getting installed apps")
			apps, err := api.GetInstalledApps()
			handleError(err, "Failed to get installed apps")

			log.Info().Msg("Getting appstores")
			appstores, err := api.GetAppstores()
			handleError(err, "Failed to get appstores")

			installedApps := make(map[string]bool)
			for _, app := range apps.Installed {
				installedApps[app.Info.Urn] = true
			}

			var dbApps []database.Apps
			db.Find(&dbApps)

			for _, dbApp := range dbApps {
				if !installedApps[dbApp.Urn] {
					log.Warn().Str("urn", dbApp.Urn).Msg("Deleting app from the database")
					db.Unscoped().Delete(&dbApp)
				}
			}

			log.Info().Msg("Comparing versions")
			appsWithUpdates := []types.App{}

			for _, app := range apps.Installed {
				// If app is up to date, ignore it
				if app.App.Version == app.Metadata.LatestVersion {
					log.Debug().Str("urn", app.Info.Urn).Msg("App is up to date, ignoring")
					continue
				}

				// If app has zeroed verions, ignore it
				if app.Metadata.LatestDockerVersion == "0.0.0" || app.Metadata.LatestVersion == 0 {
					log.Debug().Str("urn", app.Info.Urn).Msg("App has zeroed version, ignoring")
					continue
				}

				log.Debug().Interface("app", app).Msg("App has an update")

				var dbApp database.Apps
				dbRes := db.First(&dbApp, "urn = ?", app.Info.Urn)

				if dbRes.RowsAffected == 0 {
					log.Debug().Str("urn", app.Info.Urn).Msg("App not found in database, creating new entry")
					db.Create(&database.Apps{Urn: app.Info.Urn, Version: app.App.Version, LatestVersion: app.Metadata.LatestVersion})
					appsWithUpdates = append(appsWithUpdates, types.App{
						Urn:           app.Info.Urn,
						Name:          app.Info.Name,
						Version:       app.App.Version,
						DockerVersion: app.Metadata.LatestDockerVersion,
					})
				} else {
					log.Debug().Str("urn", app.Info.Urn).Msg("App found in database, checking versions")

					if dbApp.Version != app.App.Version || dbApp.LatestVersion != app.Metadata.LatestVersion {
						log.Debug().Str("urn", app.Info.Urn).Msg("Updating app in database")
						db.Model(&dbApp).Updates(database.Apps{LatestVersion: app.Metadata.LatestVersion, Version: app.App.Version})
						appsWithUpdates = append(appsWithUpdates, types.App{
							Urn:           app.Info.Urn,
							Name:          app.Info.Name,
							Version:       app.App.Version,
							DockerVersion: app.Metadata.LatestDockerVersion,
						})
					}
				}
			}

			log.Info().Msg("Sending notifications")

			for _, appWithUpdate := range appsWithUpdates {
				log.Logger.Info().Str("urn", appWithUpdate.Urn).Str("tipiVersion", strconv.Itoa(appWithUpdate.Version)).Str("dockerVersion", appWithUpdate.DockerVersion).Msg("App has an update")
				alertErr := alerts.SendAlert(&appWithUpdate, appstores.Appstores)
				handleError(alertErr, "Failed to send alert")
			}
		}

	},
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}

func init() {
	viper.SetEnvPrefix("tipimate")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	serverCmd.Flags().String("notification-url", "", "Notification URL (shoutrrr format)")
	serverCmd.Flags().String("runtipi-url", "", "Runtipi server URL")
	serverCmd.Flags().String("jwt-secret", "", "JWT secret")
	serverCmd.Flags().String("database-path", "tipimate.db", "Database path")
	serverCmd.Flags().Int("interval", 30, "Refresh interval in minutes")
	serverCmd.Flags().String("log-level", "info", "Log level (trace, debug, info, warn, error, fatal, panic)")
	serverCmd.Flags().Bool("insecure", true, "Disable TLS (https) for services like Gotify, Ntfy etc.")
	serverCmd.Flags().String("server-name", "", "Server name to use in notifications.")

	// Bind flags to viper
	viper.BindPFlags(serverCmd.Flags())

	rootCmd.AddCommand(serverCmd)
}
