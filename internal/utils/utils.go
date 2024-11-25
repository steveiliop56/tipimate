package utils

import (
	"fmt"
	"strings"
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