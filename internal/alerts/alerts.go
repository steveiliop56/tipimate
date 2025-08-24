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
		ServerName:      config.ServerName,
	}
}

type Alerts struct {
	NotificationUrl string
	RuntipiUrl      string
	Insecure        bool
	ServerName      string
}

func (alerts *Alerts) SendAlert(app *types.App, appstores []types.RuntipiAppstore) error {
	var err error

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

	service := strings.Split(alerts.NotificationUrl, "://")[0]

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

	if err != nil {
		return err
	}

	return nil
}

func (alerts *Alerts) sendDiscord(app *types.App, appstore types.RuntipiAppstore) error {
	id, _ := utils.SplitURN(app.Urn)
	appURL := fmt.Sprintf("%s/apps/%s/%s", alerts.RuntipiUrl, appstore.Slug, id)
	description := fmt.Sprintf("Your app %s from the %s appstore has an available update!\nUpdate to version `%s` (%d).", app.Name, appstore.Name, app.DockerVersion, app.Version)
	currentTime := time.Now().Format(time.RFC3339)

	var message types.DiscordMessage
	message.Embeds = []types.DiscordEmbed{
		{
			Description: description,
			Url:         appURL,
			Color:       "3126084",
			Timestamp:   currentTime,
			Footer: types.DiscordEmbedFooter{
				Text: "Updated at",
			},
		},
	}
	message.AvatarUrl = constants.RuntipiLogo
	message.Username = "Tipimate"

	if alerts.ServerName != "" {
		message.Embeds[0].Title = fmt.Sprintf("%s - %s (%s)", alerts.ServerName, app.Name, appstore.Name)
	} else {
		message.Embeds[0].Title = fmt.Sprintf("%s (%s)", app.Name, appstore.Name)
	}

	var webhook types.DiscordWebhook
	webhook.Json = true

	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?%s", alerts.NotificationUrl, queries.Encode())

	messageJson, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = shoutrrr.Send(url, string(messageJson))
	if err != nil {
		return err
	}

	return nil
}

func (alerts *Alerts) sendNtfy(app *types.App, appstore types.RuntipiAppstore) error {
	id, _ := utils.SplitURN(app.Urn)
	appURL := fmt.Sprintf("%s/apps/%s/%s", alerts.RuntipiUrl, appstore.Slug, id)
	description := fmt.Sprintf("Your app %s from the %s appstore has an available update!\nUpdate to version %s (%d).", app.Name, appstore.Name, app.DockerVersion, app.Version)

	var webhook types.NtfyWebhook
	webhook.Click = appURL

	if alerts.ServerName != "" {
		webhook.Title = fmt.Sprintf("%s - %s (%s)", alerts.ServerName, app.Name, appstore.Name)
	} else {
		webhook.Title = fmt.Sprintf("%s (%s)", app.Name, appstore.Name)
	}

	if alerts.Insecure {
		webhook.Scheme = "http"
	} else {
		webhook.Scheme = "https"
	}

	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?%s", alerts.NotificationUrl, queries.Encode())

	err = shoutrrr.Send(url, description)
	if err != nil {
		return err
	}

	return nil
}

func (alerts *Alerts) sendGotify(app *types.App, appstore types.RuntipiAppstore) error {
	id, _ := utils.SplitURN(app.Urn)
	appUrl := fmt.Sprintf("%s/apps/%s/%s", alerts.RuntipiUrl, appstore.Slug, id)
	description := fmt.Sprintf("Your app %s from the %s appstore has an available update!\nUpdate to version %s (%d).\nVisit %s for more information.", app.Name, appstore.Name, app.DockerVersion, app.Version, appUrl)

	var webhook types.GotifyWebhook
	webhook.DisableTls = alerts.Insecure

	if alerts.ServerName != "" {
		webhook.Title = fmt.Sprintf("%s - %s (%s)", alerts.ServerName, app.Name, appstore.Name)
	} else {
		webhook.Title = fmt.Sprintf("%s (%s)", app.Name, appstore.Name)
	}

	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?%s", alerts.NotificationUrl, queries.Encode())

	err = shoutrrr.Send(url, description)
	if err != nil {
		return err
	}

	return nil
}
