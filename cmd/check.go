package cmd

import (
	"fmt"
	"net/url"
	"os"
	"time"
	"tipimate/internal/api"
	"tipimate/internal/types"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Spinner
var s = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

// Check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates on your runtipi server",
	Long:  "Check for app updates on your runtipi server from your terminal",
	Run: func(cmd *cobra.Command, args []string) {
		// Spinner
		s.Suffix = "Getting apps with updates..."
		s.Start()

		// Get config
		var config types.CheckConfig
		err := viper.Unmarshal(&config)
		handleErrorSpinner(err, "Failed to parse config")

		// Validate config
		err = validator.New().Struct(config)
		handleErrorSpinner(err, "Failed to validate config")

		// Parse runtipi URL
		_, err = url.Parse(config.RuntipiUrl)
		handleErrorSpinner(err, "Invalid runtipi URL")

		// Create JWT token
		jwtToken, err := api.CreateJWT(config.JwtSecret)
		handleErrorSpinner(err, "Failed to create JWT token")

		// Get installed apps
		installedApps, err := api.GetInstalledApps(jwtToken, config.RuntipiUrl)
		handleErrorSpinner(err, "Failed to get installed apps")

		// Stop spinner
		s.Stop()

		// Check for updates
		for _, app := range installedApps.Installed {
			if app.App.Version < app.UpdateInfo.LatestVersion {
				fmt.Printf("%s Update available for %s to version %s (%d)\n", color.GreenString("↻"), app.Info.Name, app.UpdateInfo.LatestDockerVersion, app.UpdateInfo.LatestVersion)
			}
		}
	},
}

// Handle error with spinner
func handleErrorSpinner(err error, msg string) {
	if err != nil {
		s.Stop()
		fmt.Printf("%s %s\n", color.RedString("✘"), msg)
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

// Add command
func init() {
	// Viper
	viper.AutomaticEnv()

	// Flags
	checkCmd.Flags().String("runtipi", "", "The URL of your runtipi server")
	checkCmd.Flags().String("jwt-secret", "", "The JWT secret of your server")

	// Bind flags
	viper.BindEnv("runtipi", "RUNTIPI_URL")
	viper.BindEnv("jwt-secret", "JWT_SECRET")

	// Add command
	rootCmd.AddCommand(checkCmd)
}
