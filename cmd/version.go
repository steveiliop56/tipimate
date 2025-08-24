package cmd

import (
	"fmt"
	"tipimate/internal/constants"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of tipimate",
	Long:  "All software has versions. This is tipimate's",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", constants.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
