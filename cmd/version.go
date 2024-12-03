package cmd

import (
	"fmt"
	"tipimate/internal/assets"

	"github.com/spf13/cobra"
)

// Command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of TipiMate",
	Long: "All software has versions. This is TipiMate's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("TipiMate %s\n", assets.Version)
	},
}

// Add command
func init() {
	rootCmd.AddCommand(versionCmd)
}
