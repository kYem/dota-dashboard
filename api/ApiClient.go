package api

import (
	"net/http"
	"fmt"
	"github.com/kYem/dota-dashboard/config"
)

const VERSION = "v001"

type client struct {
	schema        string
	hostname      string
	game          string
	apiKey        string
}

func (client *client) GetMatchHistory(matchId string) (*http.Response) {

	url := fmt.Sprintf("%s&match_id=%s", client.getMatchUrl("GetMatchDetails"), matchId)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}



func (client *client) GetTopLiveGames(partner string) (*http.Response) {

	url := fmt.Sprintf("%s&partner=%s", client.getMatchUrl("GetTopLiveGame"), partner)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}

func (client *client) getMatchUrl(endpoint string) string {
	return fmt.Sprintf(
		"%s://%s/%s/%s/%s?key=%s",
		client.schema,
		client.hostname,
		client.game,
		endpoint,
		VERSION,
		client.apiKey,
	)
}

func (client *client) GetRealTimeStats(serverSteamId string) (*http.Response) {

	url := fmt.Sprintf("%s&server_steam_id=%s", client.getMatchStatsUrl("GetRealTimeStats"), serverSteamId)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}

func (client *client) getMatchStatsUrl(endpoint string) string {
	return fmt.Sprintf(
		"%s://%s/%s/%s/%s?key=%s",
		client.schema,
		client.hostname,
		"IDOTA2MatchStats_570",
		endpoint,
		"v1",
		client.apiKey,
	)
}


func GetClient(config config.ValveApiConfig) client {
	return client{
		config.Schema,
		config.Hostname,
		config.DotaGame,
		config.ApiKey,
	}
}
