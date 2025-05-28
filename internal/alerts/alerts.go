package alerts

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"tipimate/internal/constants"
	"tipimate/internal/types"
	"tipimate/internal/utils"

	"github.com/containrrr/shoutrrr"
	"github.com/google/go-querystring/query"
	"github.com/rs/zerolog/log"
)

func NewAlerts(config types.AlertsConfig) *Alerts {
	return &Alerts{
		NotificationUrl: config.NotificationUrl,
		RuntipiUrl:      config.RuntipiUrl,
		Insecure:        config.Insecure,
	}
}

type Alerts struct {
	NotificationUrl string
	RuntipiUrl      string
	Insecure        bool
}

func (alerts *Alerts) SendAlert(app *types.App, appstores []types.RuntipiAppstore) error {
	// Variables
	var err error

	// Get appstore
	_, slug := utils.SplitURN(app.Urn)
	appstore := utils.GetAppstore(appstores, slug)

	if appstore == nil {
		appstore = &types.RuntipiAppstore{
			Name:    "Unknown Appstore",
			Slug:    slug,
			Url:     "",
			Enabled: true,
		}
	}

	// Get notification URL service
	service := strings.Split(alerts.NotificationUrl, "://")[0]

	// Use correct service based on URL
	switch service {
	case "discord":
		log.Debug().Str("service", service).Msg("Selected Discord notification service")
		err = alerts.sendDiscord(app, *appstore)
	case "ntfy":
		log.Debug().Str("service", service).Msg("Selected Ntfy notification service")
		err = alerts.sendNtfy(app, *appstore)
	case "gotify":
		log.Debug().Str("service", service).Msg("Selected Gotify notification service")
		err = alerts.sendGotify(app, *appstore)
	default:
		log.Warn().Str("service", service).Msg("Unsupported notification service")
	}

	// Handle error
	if err != nil {
		return err
	}

	return nil
}

func (alerts *Alerts) sendDiscord(app *types.App, appstore types.RuntipiAppstore) error {
	// Variables
	id, _ := utils.SplitURN(app.Urn)
	appURL := fmt.Sprintf("%s/apps/%s/%s", alerts.RuntipiUrl, appstore.Slug, id)
	description := fmt.Sprintf("Your app %s (%s) has an available update!\nUpdate to version `%s` (%d)", app.Name, appstore.Name, app.DockerVersion, app.Version)
	currentTime := time.Now().Format(time.RFC3339)

	// Message
	var message types.DiscordMessage
	message.Embeds = []types.DiscordEmbed{
		{
			Title:       fmt.Sprintf("%s (%s)", app.Name, appstore.Name),
			Description: description,
			Url:         appURL,
			Color:       "3126084",
			Footer: types.DiscordEmbedFooter{
				Text: "Updated at",
			},
			TimeStamp: currentTime,
		},
	}
	message.AvatarUrl = constants.LogoUrl
	message.Username = "Tipimate"

	// Query params
	var webhook types.DiscordWebhook
	webhook.Json = true

	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	// Final url
	url := fmt.Sprintf("%s?%s", alerts.NotificationUrl, queries.Encode())

	// Marshal message
	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Send message
	err = shoutrrr.Send(url, string(messageJson))
	if err != nil {
		return err
	}

	return nil
}

func (alerts *Alerts) sendNtfy(app *types.App, appstore types.RuntipiAppstore) error {
	// Variables
	id, _ := utils.SplitURN(app.Urn)
	appURL := fmt.Sprintf("%s/apps/%s/%s", alerts.RuntipiUrl, appstore.Slug, id)
	description := fmt.Sprintf("Your app %s (%s) has an available update!\nUpdate to version %s (%d)", app.Name, appstore.Name, app.DockerVersion, app.Version)

	// Message
	var webhook types.NtfyWebhook
	webhook.Click = appURL
	webhook.Title = fmt.Sprintf("%s (%s)", app.Name, appstore.Name)

	if alerts.Insecure {
		webhook.Scheme = "http"
	} else {
		webhook.Scheme = "https"
	}

	// Query params
	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	// Final url
	url := fmt.Sprintf("%s?%s", alerts.RuntipiUrl, queries.Encode())

	// Send
	err = shoutrrr.Send(url, description)
	if err != nil {
		return err
	}

	return nil
}

func (alerts *Alerts) sendGotify(app *types.App, appstore types.RuntipiAppstore) error {
	// Vars
	id, _ := utils.SplitURN(app.Urn)
	appUrl := fmt.Sprintf("%s/apps/%s/%s", alerts.RuntipiUrl, appstore.Name, id)
	description := fmt.Sprintf("Your app %s (%s) has an available update!\nUpdate to version %s (%d)\nVisit %s for more information", app.Name, appstore.Name, app.DockerVersion, app.Version, appUrl)

	// Message
	var webhook types.GotifyWebhook
	webhook.Title = fmt.Sprintf("%s (%s)", app.Name, appstore.Name)
	webhook.DisableTls = alerts.Insecure

	// Query params
	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	// Final url
	url := fmt.Sprintf("%s?%s", alerts.NotificationUrl, queries.Encode())

	// Send
	err = shoutrrr.Send(url, description)
	if err != nil {
		return err
	}

	return nil
}
