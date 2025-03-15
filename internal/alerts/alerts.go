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

func SendAlert(app *types.App, notifyUrl string, runtipiUrl string, appstore string, noTls bool) error {
	// Vars
	var err error

	// Get notification URL service
	service := strings.Split(notifyUrl, "://")[0]

	// Use correct service based on URL
	switch service {
	case "discord":
		log.Debug().Str("service", service).Msg("Selected Discord notification service")
		err = SendDiscord(app, notifyUrl, runtipiUrl, appstore)
	case "ntfy":
		log.Debug().Str("service", service).Msg("Selected Ntfy notification service")
		err = SendNtfy(app, notifyUrl, runtipiUrl, appstore, noTls)
	case "gotify":
		log.Debug().Str("service", service).Msg("Selected Gotify notification service")
		err = SendGotify(app, notifyUrl, runtipiUrl, noTls)
	default:
		log.Warn().Str("service", service).Msg("Unsupported notification service")
	}

	// Handle error
	if err != nil {
		return err
	}

	return nil
}

func SendDiscord(app *types.App, discordUrl string, runtipiUrl string, appstore string) error {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", runtipiUrl, app.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version `%s` (%d)", app.Name, app.DockerVersion, app.Version)
	currentTime := time.Now().Format(time.RFC3339)

	// Message
	var message types.DiscordMessage
	message.Embeds = []types.DiscordEmbed{
		{
			Title:       app.Name,
			Description: description,
			Url:         appUrl,
			Color:       "3126084",
			Footer: types.DiscordEmbedFooter{
				Text: "Created at",
			},
			TimeStamp: currentTime,
			Thumbnail: types.DiscordEmbedThumbnail{
				Url: utils.GetAppImageUrl(app.Id, appstore),
			},
		},
	}
	message.AvatarUrl = constants.RuntipiLogoUrl
	message.Username = "Tipimate"

	// Query params
	var webhook types.DiscordWebhook
	webhook.Json = true

	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	// Final url
	url := fmt.Sprintf("%s?%s", discordUrl, queries.Encode())

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

func SendNtfy(app *types.App, ntfyUrl string, runtipiUrl string, appstore string, noTls bool) error {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", runtipiUrl, app.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version %s (%d)", app.Name, app.DockerVersion, app.Version)

	// Message
	var webhook types.NtfyWebhook
	webhook.Click = appUrl
	webhook.Icon = utils.GetAppImageUrl(app.Id, appstore)
	webhook.Title = app.Name

	if noTls {
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
	url := fmt.Sprintf("%s?%s", ntfyUrl, queries.Encode())

	// Send
	err = shoutrrr.Send(url, description)
	if err != nil {
		return err
	}

	return nil
}

func SendGotify(app *types.App, gotifyUrl string, runtipiUrl string, noTls bool) error {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", runtipiUrl, app.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version %s (%d)\nVisit %s for more information", app.Name, app.DockerVersion, app.Version, appUrl)

	// Message
	var webhook types.GotifyWebhook
	webhook.Title = app.Name
	webhook.DisableTls = noTls

	// Query params
	queries, err := query.Values(webhook)
	if err != nil {
		return err
	}

	// Final url
	url := fmt.Sprintf("%s?%s", gotifyUrl, queries.Encode())

	// Send
	err = shoutrrr.Send(url, description)
	if err != nil {
		return err
	}

	return nil
}
