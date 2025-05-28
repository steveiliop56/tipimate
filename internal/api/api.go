package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"tipimate/internal/types"

	"github.com/golang-jwt/jwt/v5"
)

func NewAPI(config types.APIConfig) (*API, error) {
	// Create JWT token
	token, err := createJWT(config.Secret)

	if err != nil {
		return nil, err
	}

	// Create transport
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure, MinVersion: tls.VersionTLS12},
	}

	// Create client
	client := http.Client{
		Transport: tr,
	}

	// Create API instance
	return &API{
		Client:     client,
		RuntipiUrl: config.RuntipiUrl,
		Token:      token,
	}, nil
}

type API struct {
	Client     http.Client
	RuntipiUrl string
	Token      string
}

func createJWT(secret string) (string, error) {
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

func (api *API) apiRequest(path string, method string) (*http.Response, error) {
	// Define variables
	url := fmt.Sprintf("%s%s", api.RuntipiUrl, path)
	bearer := fmt.Sprintf("Bearer %s", api.Token)

	// Create request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Authorization", bearer)

	// Send request
	res, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Check status code
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status code: %d", res.StatusCode)
	}

	// Return response
	return res, nil
}

func (api *API) GetInstalledApps() (types.GetInstalledAppsResponse, error) {
	// Create response var
	var installedApps types.GetInstalledAppsResponse

	// Get response from API
	res, err := api.apiRequest("/api/apps/installed", "GET")

	if err != nil {
		return installedApps, err
	}

	defer res.Body.Close()

	// Decode response
	err = json.NewDecoder(res.Body).Decode(&installedApps)
	if err != nil {
		return installedApps, err
	}

	// Return response
	return installedApps, nil
}

func (api *API) GetAppstores() (types.GetAppstoresResponse, error) {
	// Create response var
	var appstores types.GetAppstoresResponse

	// Get response from API
	res, err := api.apiRequest("/api/marketplace/enabled", "GET")

	if err != nil {
		return appstores, err
	}

	defer res.Body.Close()

	// Decode response
	err = json.NewDecoder(res.Body).Decode(&appstores)
	if err != nil {
		return appstores, err
	}

	// Return response
	return appstores, nil
}
