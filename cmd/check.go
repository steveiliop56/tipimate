package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"tipimate/internal/api"
	"tipimate/internal/assets"
	"tipimate/internal/spinner"
	"tipimate/internal/types"
	"tipimate/internal/utils"

	"github.com/fatih/color"
	"github.com/gookit/validate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Viper
var checkViper = viper.New()

// Command
var checkCmd = &cobra.Command{
	Use:  "check",
	Short: "Check for updates on your Runtipi server",
	Long: "Check for app updates on your Runtipi server from your terminal",
	Run: func(cmd *cobra.Command, args []string) {
		// Spinner
		spinner.Init()
		spinner.Start()
		spinner.Info(fmt.Sprintf("Tipicord version %s", assets.Version))
		spinner.SetMessage("Getting apps with updates...")

		// Get config
		var config types.CheckConfig
		viperParseErr := checkViper.Unmarshal(&config)
		utils.HandleErrorSpinner("Failed to parse config", viperParseErr)

		// Validate config
		validtor := validate.Struct(config)
		if !validtor.Validate() {
			utils.HandleErrorSpinner("Invalid config", errors.New(validtor.Errors.One()))
		}

		// Parse runtipi URL
		_, runtipiParseErr := url.Parse(config.RuntipiUrl)
		utils.HandleErrorSpinner("Invalid Runtipi server URL", runtipiParseErr)

		// Create JWT token
		jwtToken, jwtErr := api.CreateJWT(config.JwtSecret)
		utils.HandleErrorSpinner("Failed to create JWT token", jwtErr)

		// Get installed apps
		installedApps, installedAppsErr := api.GetInstalledApps(jwtToken, config.RuntipiUrl)
		utils.HandleErrorSpinner("Failed to get installed apps", installedAppsErr)

		// Check for updates
		for _, app := range installedApps.Installed {
			if app.App.Version < app.UpdateInfo.LatestVersion {
				spinner.Custom(fmt.Sprintf("Update available for %s to version %s (%d)", app.Info.Name, app.UpdateInfo.LatestDockerVersion, app.UpdateInfo.LatestVersion), color.GreenString("↻"))
			}
		}

		// Stop spinner
		spinner.Stop()

	},
}

// Add command
func init() {
	checkViper.AutomaticEnv()
	checkCmd.Flags().String("runtipi", "", "Runtipi server URL")
	checkCmd.Flags().String("jwt-secret", "", "Runtipi JWT secret")
	checkViper.BindEnv("jwt-secret", "JWT_SECRET")
	checkViper.BindPFlags(checkCmd.Flags())
	rootCmd.AddCommand(checkCmd)
}
