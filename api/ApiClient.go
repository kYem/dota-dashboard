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

func (this *client) GetMatchHistory(matchId string) (*http.Response) {

	url := fmt.Sprintf("%s&match_id=%s", this.getMatchUrl("GetMatchDetails"), matchId)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}



func (this *client) GetTopLiveGames(partner string) (*http.Response) {

	url := fmt.Sprintf("%s&partner=%s", this.getMatchUrl("GetTopLiveGame"), partner)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}

func (this *client) getMatchUrl(endpoint string) string {
	return fmt.Sprintf(
		"%s://%s/%s/%s/%s?key=%s",
		this.schema,
		this.hostname,
		this.game,
		endpoint,
		VERSION,
		this.apiKey,
	)
}

func (this *client) GetRealTimeStats(serverSteamId string) (*http.Response) {

	url := fmt.Sprintf("%s&server_steam_id=%s", this.getMatchStatsUrl("GetRealTimeStats"), serverSteamId)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}

func (this *client) getMatchStatsUrl(endpoint string) string {
	return fmt.Sprintf(
		"%s://%s/%s/%s/%s?key=%s",
		this.schema,
		this.hostname,
		"IDOTA2MatchStats_570",
		endpoint,
		"v1",
		this.apiKey,
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