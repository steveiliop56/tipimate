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
	token, err := createJWT(config.Secret)

	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: config.Insecure, MinVersion: tls.VersionTLS12},
	}

	client := http.Client{
		Transport: tr,
	}

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "cli",
	})

	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signed, err
}

func (api *API) apiRequest(path string, method string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", api.RuntipiUrl, path)
	bearer := fmt.Sprintf("Bearer %s", api.Token)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", bearer)

	res, err := api.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status code: %d", res.StatusCode)
	}

	return res, nil
}

func (api *API) GetInstalledApps() (types.GetInstalledAppsResponse, error) {
	var installedApps types.GetInstalledAppsResponse

	res, err := api.apiRequest("/api/apps/installed", "GET")

	if err != nil {
		return installedApps, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&installedApps)
	if err != nil {
		return installedApps, err
	}

	return installedApps, nil
}

func (api *API) GetAppstores() (types.GetAppstoresResponse, error) {
	var appstores types.GetAppstoresResponse

	res, err := api.apiRequest("/api/marketplace/enabled", "GET")

	if err != nil {
		return appstores, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&appstores)
	if err != nil {
		return appstores, err
	}

	return appstores, nil
}
