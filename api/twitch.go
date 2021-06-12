package api

import (
	"errors"
	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
	"os"
)

type Twitch struct {
	client *helix.Client
}

var TwitchClient *Twitch

func init() {
	TwitchClient = &Twitch{client: createTwitchClient()}
}

const DotaGameId = "29595"

func (t *Twitch) GetStreams(userIds []string, first int) (*helix.StreamsResponse, error) {
	params := &helix.StreamsParams{
		Language: []string{"en"},
		GameIDs:  []string{DotaGameId},
		UserIDs: userIds,
	}

	if first > 0 {
		params.First = first
	}

	streams, err := t.client.GetStreams(params)

	if err == nil && streams.Error != "" {
		return streams, errors.New("Twitch error: " + streams.ErrorMessage)
	}

	return streams, err
}

func createTwitchClient() *helix.Client {
	client, err := helix.NewClient(&helix.Options{
		ClientID: os.Getenv("TWITCH_API_KEY"),
		ClientSecret: os.Getenv("TWITCH_API_SECRET"),
	})
	if err != nil {
		log.Fatal(err)
	}

	var scopes []string
	token, errToken := client.RequestAppAccessToken(scopes)
	if errToken != nil {
		log.Error(err)
	}
	client.SetAppAccessToken(token.Data.AccessToken)

	return client
}


