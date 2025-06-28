package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
	"tipimate/internal/api"
	"tipimate/internal/types"
	"tipimate/internal/utils"

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
		s.Suffix = " Getting apps with updates..."
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

		// Create API client
		apiConfig := types.APIConfig{
			RuntipiUrl: config.RuntipiUrl,
			Secret:     config.JwtSecret,
			Insecure:   config.Insecure,
		}

		api, err := api.NewAPI(apiConfig)
		handleErrorSpinner(err, "Failed to create API client")

		// Get installed apps
		apps, err := api.GetInstalledApps()
		handleErrorSpinner(err, "Failed to get installed apps")

		// Get appstores
		appstores, err := api.GetAppstores()
		handleErrorSpinner(err, "Failed to get appstores")

		// Stop spinner
		s.Stop()

		updatesAvailable := false

		// Check for updates
		for _, app := range apps.Installed {
			// If app has zeroed verions, ignore it
			if app.Metadata.LatestDockerVersion == "0.0.0" || app.Metadata.LatestVersion == 0 {
				continue
			}

			// Check if app is up to date
			if app.App.Version < app.Metadata.LatestVersion {
				updatesAvailable = true
				_, slug := utils.SplitURN(app.Info.Urn)
				appstore := utils.GetAppstore(appstores.Appstores, slug)
				if appstore == nil {
					fmt.Printf("%s Update available for the app %s from the %s appstore to version %s (%d)!\n", color.GreenString("↻"), app.Info.Name, "Unknown Appstore", app.Metadata.LatestDockerVersion, app.Metadata.LatestVersion)
					continue
				}
				fmt.Printf("%s Update available for the app %s from the %s appstore to version %s (%d)!\n", color.GreenString("↻"), app.Info.Name, appstore.Name, app.Metadata.LatestDockerVersion, app.Metadata.LatestVersion)
			}
		}

		if !updatesAvailable {
			fmt.Printf("%s All apps are up to date!\n", color.GreenString("✔"))
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
	viper.SetEnvPrefix("tipimate")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Flags
	checkCmd.Flags().String("runtipi-url", "", "Runtipi server URL")
	checkCmd.Flags().String("jwt-secret", "", "JWT secret")
	checkCmd.Flags().Bool("insecure", false, "Ignore self-signed certificates")

	// Bind flags to viper
	viper.BindPFlags(checkCmd.Flags())

	// Add command
	rootCmd.AddCommand(checkCmd)
}
