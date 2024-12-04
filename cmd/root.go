package cmd

import (
	"fmt"
	"os"

	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
)

// Main command
var rootCmd = &cobra.Command{
	Use:   "tipimate",
	Short: "Discord notifications for your runtipi server",
	Long: "TipiMate is a simple tool that monitors your runtipi server for app updates and notifies you via Discord notifications",
}

// Execute command
func Execute() {
	cc.Init(&cc.Config{
        RootCmd:       rootCmd,
        Headings:      cc.Blue,
        Example:       cc.Italic,
		Commands: 	   cc.Green,
        ExecName:      cc.White,
        Flags:         cc.White,
    })
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("An error occured while executing, error: %s\n", err.Error())
		os.Exit(1)
	}
}
