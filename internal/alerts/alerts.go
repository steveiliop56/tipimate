package alerts

import (
	"encoding/json"
	"fmt"
	"time"
	"tipimate/internal/constants"
	"tipimate/internal/types"
	"tipimate/internal/utils"

	"github.com/containrrr/shoutrrr"
	"github.com/google/go-querystring/query"
)

func SendDiscord(app *types.SimpleApp, discordUrl string, runtipiUrl string, appstore string) (error) {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", runtipiUrl, app.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version `%s` (%d)", app.Name, app.DockerVersion, app.Version)
	currentTime := time.Now().Format(time.RFC3339)

	// Message
	var message types.Message
	message.Embeds = []types.Embed{
		{
			Title: app.Name,
			Description: description,
			Url: appUrl,
			Color: "3126084",
			Footer: types.EmbedFooter{
				Text: "Created at",
			},
			TimeStamp: currentTime,
			Thumbnail: types.EmbedThumbnail{
				Url: utils.GetAppImageUrl(app.Id, appstore),
			},
		},
	}
	message.AvatarUrl = constants.RuntipiLogoUrl
	message.Username = "TipiMate"

	// Query params
	var webhook types.Webhook
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
