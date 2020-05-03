package api

import (
	"encoding/json"
	"fmt"
	"github.com/kYem/dota-dashboard/config"
	"github.com/kYem/dota-dashboard/dota"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const VERSION = "v001"

type SteamClient struct {
	schema        string
	hostname      string
	game          string
	apiKey        string
}


var SteamApi SteamClient

func init() {
	SteamApi = GetClient(config.LoadConfig())
}

func (client *SteamClient) GetMatchHistory(matchId string) *http.Response {

	url := fmt.Sprintf("%s&match_id=%s", client.getMatchUrl("GetMatchDetails"), matchId)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("GetMatchHistory: failed: %s", err)
	}

	return resp
}



// 0 NA or combined?
// 1 China
// 2 Europe (West)
// 3 South America
func (client *SteamClient) GetTopLiveGames(partner string) *http.Response {

	url := fmt.Sprintf("%s&partner=%s", client.getMatchUrl("GetTopLiveGame"), partner)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("GetTopLiveGames: failed: %s", err)
	}

	return resp
}

func (client *SteamClient) getMatchUrl(endpoint string) string {
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

func (client *SteamClient) GetRealTimeStats(serverSteamId string) *http.Response {

	url := fmt.Sprintf("%s&server_steam_id=%s", client.getMatchStatsUrl("GetRealTimeStats"), serverSteamId)
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	return resp
}

func (client *SteamClient) GetLiveLeagueGames() *http.Response {
	resp, err := http.Get(client.getMatchUrl("GetLiveLeagueGames"))
	if err != nil {
		panic(err)
	}

	return resp
}

func (client *SteamClient) getMatchStatsUrl(endpoint string) string {
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

func (client *SteamClient) getEconUrl(endpoint string) string {
	return fmt.Sprintf(
		"%s://%s/%s/%s/%s?key=%s",
		client.schema,
		client.hostname,
		"IEconDOTA2_570",
		endpoint,
		"v1",
		client.apiKey,
	)
}

func (client *SteamClient) GetHeroes() ([]dota.HeroBasic, error) {
	resp, err := http.Get(client.getEconUrl("GetHeroes"))


	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, err
	}

	var heroes dota.GetHeroes
	if err := json.NewDecoder(resp.Body).Decode(&heroes); err != nil {
		log.Println(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
	}

	return heroes.Result.Heroes, err
}

func GetClient(config config.ValveApiConfig) SteamClient {
	return SteamClient{
		config.Schema,
		config.Hostname,
		config.DotaGame,
		config.ApiKey,
	}
}
