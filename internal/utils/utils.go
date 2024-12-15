package utils

import (
	"fmt"
	"os"
	"strings"
	"tipimate/internal/spinner"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func GetAppImageUrl(appId string, appstore string) (string) {
	// Vars
	branch := "master"

	// Check if other branch is used
	if strings.Contains(appstore, "tree") {
		branch = strings.Split(appstore, "tree/")[1]
	}

	// Return image
	return fmt.Sprintf("%s/blob/%s/apps/%s/metadata/logo.jpg?raw=true", appstore, branch, appId)
}

func HandleErrorSpinner(msg string, err error) {
	// Check error
	if err != nil {
		// Stop spinner
		spinner.Fail(msg)
		// Print error
		fmt.Printf("Error: %s\n", err)
		// Exit
		os.Exit(1)
	}
}

func HandleErrorLogger(err error, msg string) {
	// Check error
	if err != nil {
		// Log error
		log.Error().Str("err", err.Error()).Msg(msg)
		// Exit
		os.Exit(1)
	}
}

func HandleErrorLoggerNoExit(err error, msg string) {
	// Check error
	if err != nil {
		// Log error
		log.Error().Str("err", err.Error()).Msg(msg)
	}
}

func HandleError(err error, msg string) {
	// Check error
	if err != nil {
		// Print error
		fmt.Printf("Error: %s\n", err)
		// Exit
		os.Exit(1)
	}
}

func GetLogLevel(level string) (zerolog.Level) {
	// Get level from string
	switch level {
		case "panic":
			return zerolog.PanicLevel
		case "fatal":
			return zerolog.FatalLevel
		case "error":
			return zerolog.ErrorLevel
		case "warn":
			return zerolog.WarnLevel
		case "info":
			return zerolog.InfoLevel
		case "debug":
			return zerolog.DebugLevel
		case "trace":
			return zerolog.TraceLevel
	}

	// Fallback to info
	return zerolog.InfoLevel
}