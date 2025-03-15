package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"tipimate/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret string) (string, error) {
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "cli",
	})

	// Sign token with secret
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	// Return token
	return signed, err
}

func GetInstalledApps(token string, runtipiUrl string) (types.GetInstalledAppsResponse, error) {
	// Define vars
	var appsUrl = fmt.Sprintf("%s/api/apps/installed", runtipiUrl)
	var bearer = fmt.Sprintf("Bearer %s", token)

	// Create response var
	var response types.GetInstalledAppsResponse

	// Create transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Create client
	client := http.Client{
		Transport: tr,
	}

	// Create request
	req, err := http.NewRequest("GET", appsUrl, nil)
	if err != nil {
		return response, err
	}

	// Set headers
	req.Header.Set("Authorization", bearer)

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer res.Body.Close()

	// Decode response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return response, err
	}

	// Return response
	return response, nil
}
