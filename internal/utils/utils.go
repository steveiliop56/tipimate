package utils

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
)

func GetAppImageUrl(appId string, appstore string) string {
	// Vars
	branch := "master"

	// Check if other branch is used
	if strings.Contains(appstore, "tree") {
		branch = strings.Split(appstore, "tree/")[1]
	}

	// Return image
	return fmt.Sprintf("%s/blob/%s/apps/%s/metadata/logo.jpg?raw=true", appstore, branch, appId)
}

func GetLogLevel(level string) zerolog.Level {
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
	default:
		return zerolog.InfoLevel
	}
}
