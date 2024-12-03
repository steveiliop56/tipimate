package cmd

import (
	"fmt"
	"tipicord/internal/assets"

	"github.com/spf13/cobra"
)

// Command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of TipiCord",
	Long: "All software has versions. This is TipiCord's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TipiCord %s\n", assets.Version)
	},
}

// Add command
func init() {
	rootCmd.AddCommand(versionCmd)
}
