package alerts

import (
	"encoding/json"
	"fmt"
	"time"
	"tipicord/internal/constants"
	"tipicord/internal/types"
	"tipicord/internal/utils"

	"github.com/containrrr/shoutrrr"
	"github.com/google/go-querystring/query"
)

func SendAppUpdateAlert(info *types.AppUpdateAlert) (error) {
	// Vars
	appUrl := fmt.Sprintf("%s/apps/%s", info.ServerUrl, info.Id)
	description := fmt.Sprintf("Your app %s has an available update!\nUpdate to version `%s` (%d)", info.Name, info.DockerVersion, info.Version)
	currentTime := time.Now().Format(time.RFC3339)

	// Message
	var message types.Message
	message.Embeds = []types.Embed{
		{
			Title: info.Name,
			Description: description,
			Url: appUrl,
			Color: "3126084",
			Footer: types.EmbedFooter{
				Text: "Created at",
			},
			TimeStamp: currentTime,
			Thumbnail: types.EmbedThumbnail{
				Url: utils.GetAppImageUrl(info.Id, info.AppStore),
			},
		},
	}
	message.AvatarUrl = constants.RuntipiLogoUrl
	message.Username = "TipiCord"

	// Query params
	var webhook types.Webhook
	webhook.Json = true
	
	queries, queriesErr := query.Values(webhook)
	if queriesErr != nil {
		return queriesErr
	}

	// Final url
	url := fmt.Sprintf("%s?%s", info.DiscordUrl, queries.Encode())

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
