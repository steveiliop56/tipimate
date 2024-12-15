package cmd

import (
	"errors"
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
	"github.com/gookit/validate"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Viper
var serverViper = viper.New()

// Command
var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start the Tipicord server",
	Long: "Use the server command to automatically check for updates on your Runtipi server and send them to your Discord server",
	Run: func(cmd *cobra.Command, args []string) {
		// Logger
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger().Level(zerolog.DebugLevel)
		log.Info().Str("version", assets.Version).Msg("Starting Tipicord server")
		
		// Get config
		var config types.ServerConfig
		viperParseErr := serverViper.Unmarshal(&config)
		utils.HandleErrorLogger(viperParseErr, "Failed to parse config")
		if config.RuntipiInternalUrl == "" {
			config.RuntipiInternalUrl = config.RuntipiUrl
		}

		// Validate config
		validtor := validate.Struct(config)
		if !validtor.Validate() {
			utils.HandleErrorLogger(errors.New(validtor.Errors.One()), "Invalid config")
		}

		// Validate URL
		sr := router.ServiceRouter{}
		_, urlValidateErr := sr.Locate(config.NotifyUrl)
		utils.HandleErrorLogger(urlValidateErr, "Invalid notification URL")

		// Validate runtipi URL
		_, runtipiParseErr := url.Parse(config.RuntipiUrl)
		utils.HandleErrorLogger(runtipiParseErr, "Invalid Runtipi server URL")

		// Initialize db
		db, dbErr := database.InitDb(config.DbPath)
		utils.HandleErrorLogger(dbErr, "Failed to initialize database")

		// Create JWT
		token, jwtErr := api.CreateJWT(config.JwtSecret)
		utils.HandleErrorLogger(jwtErr, "Failed to create JWT token")

		// Main loops
		for {
			log.Info().Msg("Checking for updates")

			// Get installed apps
			log.Info().Msg("Getting installed apps")
			installedApps, installedAppsErr := api.GetInstalledApps(token, config.RuntipiInternalUrl)
			if installedAppsErr != nil {
				utils.HandleErrorLoggerNoExit(installedAppsErr, "Failed to get installed apps")
				continue
			}

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
			appsWithUpdates := []types.SimpleApp{}

			for _, app := range installedApps.Installed {
				// Get app from database
				var dbApp database.Schema
		 		dbRes := db.First(&dbApp, "id = ?", app.App.Id)

				// Check if app is not in database
				if dbRes.RowsAffected == 0 {
					// Create app in database
					db.Create(&database.Schema{ Id: app.App.Id, Version: app.App.Version, LatestVersion: app.UpdateInfo.LatestVersion })

					// Check for updates
					if app.App.Version < app.UpdateInfo.LatestVersion {
						appsWithUpdates = append(appsWithUpdates, types.SimpleApp{
							Id: app.App.Id,
							Name: app.Info.Name,
							Version: app.App.Version,
							DockerVersion: app.UpdateInfo.LatestDockerVersion,
						})
					}
				} else {
					// Check if app has changed
					if dbApp.LatestVersion != app.UpdateInfo.LatestVersion || dbApp.Version != app.App.Version {
						// Check for updates
						if app.App.Version < app.UpdateInfo.LatestVersion {
							appsWithUpdates = append(appsWithUpdates, types.SimpleApp{
								Id: app.App.Id,
								Name: app.Info.Name,
								Version: app.App.Version,
								DockerVersion: app.UpdateInfo.LatestDockerVersion,
							})
						}

						// Modify db
						db.Model(&dbApp).Updates(database.Schema{ LatestVersion: app.UpdateInfo.LatestVersion, Version: app.App.Version })
					}
				}
			}

			// Send notifications
			log.Info().Msg("Sending notifications")
			for _, appWithUpdate := range appsWithUpdates {
				// Log
				log.Logger.Info().Str("appId", appWithUpdate.Id).Str("tipiVersion", strconv.Itoa(appWithUpdate.Version)).Str("dockerVersion", appWithUpdate.DockerVersion).Msg("App has an update")

				// Send alert
				alertErr := alerts.SendAlert(&appWithUpdate, config.NotifyUrl, config.RuntipiUrl, config.Appstore)

				// Handle error
				utils.HandleErrorLoggerNoExit(alertErr, "Failed to send app update alert")
			}

			// Sleep
			time.Sleep(time.Duration(config.Refresh) * time.Minute)
		}

	},
}

// Add command
func init() {
	serverViper.AutomaticEnv()
	serverCmd.Flags().String("notify-url", "", "Notification URL (shoutrrr format)")
	serverCmd.Flags().String("runtipi", "", "Runtipi server URL")
	serverCmd.Flags().String("runtipi-internal", "", "Runtipi internal URL (used when running in the same server as runtipi)")
	serverCmd.Flags().String("jwt-secret", "", "JWT secret")
	serverCmd.Flags().String("appstore", "https://github.com/runtipi/runtipi-appstore", "Appstore URL for images")
	serverCmd.Flags().String("db-path", "tipimate.db", "Database path")
	serverCmd.Flags().Int("refresh", 30, "Refresh interval")
	serverViper.BindEnv("notify-url", "NOTIFY_URL")
	serverViper.BindEnv("jwt-secret", "JWT_SECRET")
	serverViper.BindEnv("db-path", "DB_PATH")
	serverViper.BindEnv("runtipi-internal", "RUNTIPI_INTERNAL")
	serverViper.BindPFlags(serverCmd.Flags())
	rootCmd.AddCommand(serverCmd)
}
