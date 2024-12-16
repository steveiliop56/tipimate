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

func SendAlert(app *types.SimpleApp, notifyUrl string, runtipiUrl string, appstore string, noTls bool) (error) {
	// Vars
	var sendErr error

	// Get notification URL service
	service := strings.Split(notifyUrl, "://")[0]

	// Use correct service based on URL
	switch service {
		case "discord":
			log.Debug().Str("service", service).Msg("Selected Discord notification service")
			sendErr = SendDiscord(app, notifyUrl, runtipiUrl, appstore)
		case "ntfy":
			log.Debug().Str("service", service).Msg("Selected Ntfy notification service")
			scheme := "https"
			if noTls {
				scheme = "http"
			}
			sendErr = SendNtfy(app, notifyUrl, runtipiUrl, appstore, scheme)
		default:
			log.Warn().Str("service", service).Msg("Unsupported notification service")
	}

	// Handle error
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func SendDiscord(app *types.SimpleApp, discordUrl string, runtipiUrl string, appstore string) (error) {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", runtipiUrl, app.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version `%s` (%d)", app.Name, app.DockerVersion, app.Version)
	currentTime := time.Now().Format(time.RFC3339)

	// Message
	var message types.DiscordMessage
	message.Embeds = []types.DiscordEmbed{
		{
			Title: app.Name,
			Description: description,
			Url: appUrl,
			Color: "3126084",
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
	message.Username = "TipiMate"

	// Query params
	var webhook types.DiscordWebhook
	webhook.Json = true
	
	queries, queriesErr := query.Values(webhook)
	if queriesErr != nil {
		return queriesErr
	}

	// Final url
	url := fmt.Sprintf("%s?%s", discordUrl, queries.Encode())

	// Marshal message
	messageJson, marshalErr := json.Marshal(message)
	if marshalErr != nil {
		return marshalErr
	}

	// Send message
	sendErr := shoutrrr.Send(url, string(messageJson))
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func SendNtfy(app *types.SimpleApp, ntfyUrl string, runtipiUrl string, appstore string, scheme string) (error) {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", runtipiUrl, app.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version %s (%d)", app.Name, app.DockerVersion, app.Version)

	// Message
	var webhook types.NtfyWebhook
	webhook.Click = appUrl
	webhook.Icon = utils.GetAppImageUrl(app.Id, appstore)
	webhook.Title = app.Name
	webhook.Scheme = scheme

	// Query params
	queries, queriesErr := query.Values(webhook)
	if queriesErr != nil {
		return queriesErr
	}

	// Final url
	url := fmt.Sprintf("%s?%s", ntfyUrl, queries.Encode())

	// Send
	sendErr := shoutrrr.Send(url, description)
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func SendGotify(app *types.SimpleApp, gotifyUrl string, runtipiUrl string, noTls bool) (error) {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", runtipiUrl, app.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version %s (%d)\nVisit %s for more information", app.Name, app.DockerVersion, app.Version, appUrl)

	// Message
	var webhook types.GotifyWebhook
	webhook.Title = app.Name
	webhook.DisableTls = noTls

	// Query params
	queries, queriesErr := query.Values(webhook)
	if queriesErr != nil {
		return queriesErr
	}

	// Final url
	url := fmt.Sprintf("%s?%s", gotifyUrl, queries.Encode())

	// Send
	sendErr := shoutrrr.Send(url, description)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
