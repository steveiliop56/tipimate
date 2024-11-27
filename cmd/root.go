package cmd

import (
	"fmt"
	"net/url"
	"os"
	"time"
	"tipicord/internal/alerts"
	"tipicord/internal/api"
	"tipicord/internal/assets"
	"tipicord/internal/database"
	"tipicord/internal/types"

	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/gookit/validate"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Viper
var cmdViper = viper.New()

// Logger
var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).With().Timestamp().Logger().Level(zerolog.DebugLevel)

// Main command
var rootCmd = &cobra.Command{
	Use:   "tipicord",
	Short: "Discord notifications for your runtipi server",
	Long: `TipiCord is a simple tool that monitors your runtipi server for app updates and notifies you via Discord notifications.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Start message
		logger.Info().Msg("Starting")

		// Get config
		logger.Info().Msg("Getting config")
		var config types.Config
		viperParseErr := cmdViper.Unmarshal(&config)
		handleError(viperParseErr, "Failed to parse config", true)

		// Validate config
		logger.Info().Msg("Validating config")
		validtor := validate.Struct(config)
		if !validtor.Validate() {
			logger.Error().Str("err", validtor.Errors.One()).Msg("Invalid config")
			os.Exit(1)
		}

		// Print version
		logger.Info().Msg(fmt.Sprintf("TipiCord %s", assets.Version))
		
		// Validate discord URL
		logger.Info().Msg("Validating Discord webhook URL")
		sr := router.ServiceRouter{}
		_, discordValidateErr := sr.Locate(config.DiscordUrl)
		handleError(discordValidateErr, "Invalid Discord webhook URL", true)

		// Parse runtipi URL
		logger.Info().Msg("Parsing Runtipi server URL")
		_, runtipiParseErr := url.Parse(config.RuntipiUrl)
		handleError(runtipiParseErr, "Invalid Runtipi server URL", true)

		// Initialize db
		logger.Info().Msg("Initializing database")
		db, dbErr := database.InitDb(config.DbPath)
		handleError(dbErr, "Failed to initialize database", true)

		// Create JWT token
		logger.Info().Msg("Creating JWT token")
		jwtToken, jwtErr := api.CreateJWT(config.JwtSecret)
		handleError(jwtErr, "Failed to create JWT token", true)

		// Main loop
		for {
			logger.Info().Msg("Checking for app updates")

			// Get installed apps
			logger.Info().Msg("Getting installed apps")
			installedApps, installedAppsErr := api.GetInstalledApps(jwtToken, config.RuntipiUrl)
			handleError(installedAppsErr, "Failed to get installed apps", true)

			// Purge uninstalled apps from db
			logger.Info().Msg("Purging uninstalled apps from db")

			installedAppIds := make(map[string]bool)
			for _, app := range installedApps.Installed {
				installedAppIds[app.App.Id] = true
			}

			var dbApps []database.Schema
			db.Find(&dbApps)

			for _, dbApp := range dbApps {
				if !installedAppIds[dbApp.Id] {
					logger.Info().Str("appId", dbApp.Id).Msg("Deleting app from the database")
					db.Unscoped().Delete(&dbApp)
				}
			}

			// Create map with apps that need to be updated
			logger.Info().Msg("Getting apps with updates")

			var appsWithUpdate []types.SimpleApp

			for _, app := range installedApps.Installed {
				var dbApp database.Schema
				dbRes := db.First(&dbApp, "id = ?", app.App.Id)

				if dbRes.RowsAffected == 0 {
					db.Create(&database.Schema{ Id: app.App.Id, Version: app.App.Version, LatestVersion: app.UpdateInfo.LatestVersion })
					if app.UpdateInfo.LatestVersion > app.App.Version {
						appsWithUpdate = append(appsWithUpdate, types.SimpleApp{
							Id: app.App.Id,
							Name: app.Info.Name,
							DockerVersion: app.UpdateInfo.LatestDockerVersion,
							Version: app.App.Version,
						})
					}
				} else {
					if dbApp.LatestVersion != app.UpdateInfo.LatestVersion || dbApp.Version != app.App.Version {
						if app.UpdateInfo.LatestVersion > app.App.Version {
							appsWithUpdate = append(appsWithUpdate, types.SimpleApp{
								Id: app.App.Id,
								Name: app.Info.Name,
								DockerVersion: app.UpdateInfo.LatestDockerVersion,
								Version: app.App.Version,
							})
						}

						db.Model(&dbApp).Updates(database.Schema{ LatestVersion: app.UpdateInfo.LatestVersion, Version: app.App.Version })
					}
				}
			}

			// Send notifications for apps with updates
			logger.Info().Msg("Sending notifications for apps with updates")
			for _, app := range appsWithUpdate {
				alertErr := alerts.SendAppUpdateAlert(&types.AppUpdateAlert{
					Id: app.Id,
					Name: app.Name,
					DockerVersion: app.DockerVersion,
					Version: app.Version,
					ServerUrl: config.RuntipiUrl,
					DiscordUrl: config.DiscordUrl,
					AppStore: config.Appstore,
				})
				handleError(alertErr, "Failed to send app update alert", false)
			}

			// Sleep
			time.Sleep(time.Duration(config.Refresh) * time.Minute)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("An error occured while executing, error: %s\n", err.Error())
		os.Exit(1)
	}
}

func handleError(err error, msg string, exit bool) {
	if err != nil {
		logger.Error().Str("err", err.Error()).Msg(msg)
		if exit {
			os.Exit(1)
		}
	}
}

func init() {
	cmdViper.AutomaticEnv()
	rootCmd.Flags().String("discord", "", "Discord webhook URL")
	rootCmd.Flags().String("runtipi", "", "Runtipi server URL")
	rootCmd.Flags().String("jwtSecret", "", "JWT secret")
	rootCmd.Flags().String("appstore", "https://github.com/runtipi/runtipi-appstore", "Runtipi appstore URL (default https://github.com/runtipi/runtipi-appstore)")
	rootCmd.Flags().String("databasePath", "tipicord.db", "Database path (default tipicord.db)")
	rootCmd.Flags().IntP("refresh", "r", 30, "Refresh interval in minutes (default 30)")
	cmdViper.BindEnv("jwtSecret", "JWT_SECRET")
	cmdViper.BindEnv("databasePath", "DATABASE_PATH")
	cmdViper.BindPFlags(rootCmd.Flags())
}


