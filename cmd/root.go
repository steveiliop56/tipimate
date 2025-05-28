package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Main command
var rootCmd = &cobra.Command{
	Use:   "tipimate",
	Short: "App update notifications for your runtipi server",
	Long:  "Tipimate is a simple tool that sends you notification when your runtipi apps have an available update",
}

// Execute command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("An error occured while executing, error: %s\n", err.Error())
		os.Exit(1)
	}
}
