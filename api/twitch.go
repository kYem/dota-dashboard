package api

import (
	"github.com/nicklaw5/helix"
	"log"
	"os"
)

type Client struct {
	apiKey        string
}

const DotaGameId = "29595"

func CreateTwitchClient() *helix.Client {
	apiClient, err := helix.NewClient(&helix.Options{
		ClientID: os.Getenv("TWITCH_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	return apiClient

}


