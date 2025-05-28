package utils

import (
	"strings"
	"tipimate/internal/types"

	"github.com/rs/zerolog"
)

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

func SplitURN(urn string) (string, string) {
	// Split URN by colon
	parts := strings.Split(urn, ":")
	if len(parts) != 2 {
		return "", ""
	}

	// Return parts
	return parts[0], parts[1]
}

func GetAppstore(appstores []types.RuntipiAppstore, slug string) *types.RuntipiAppstore {
	// Loop through appstores
	for _, appstore := range appstores {
		// Check if slug matches
		if appstore.Slug == slug {
			return &appstore
		}
	}

	// Return nil if not found
	return nil
}
