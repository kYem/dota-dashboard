package stream

import (
	"encoding/json"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/dota"
	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

const (
	gorgc         = 56939869
	cancel        = 141690233
	Kingteka      = 242697694
	SexyBamboe    = 20321748
	PurgeGamers   = 66296404
	absolutHabibi = 140873889
	JustCooman    = 175463659
	SyndereN      = 4281729
	KheZu         = 169025618
	Ghostik       = 82297531
	Mason         = 315657960
	Funkefal      = 101117815
	GunnarDotA2   = 126238768
	bobruhatv     = 86953944
	Mickee        = 152962063
	MickeeTwo     = 106755427
	Febby         = 112377459
	Mage		  = 178366364
)

var dotaToTwitchMap = map[int]string{
	gorgc:         "108268890",
	cancel:        "83195409",
	Kingteka:      "127007669",
	SexyBamboe:    "22580017",
	PurgeGamers:   "22561231",
	absolutHabibi: "140873889",
	JustCooman:    "63667409",
	SyndereN:      "26656197",
	KheZu:         "25199180",
	Ghostik:       "25199180",
	Mason:         "40754777",
	Funkefal:      "126104914",
	GunnarDotA2:   "131202285",
	bobruhatv:     "116741333",
	Mickee:        "266316098",
	MickeeTwo:     "266316098",
	Febby:         "87822995",
	Mage:          "85002144",
}

var reverseLookup = map[string]int{}

var proPlayers []dota.ProPlayer

func init() {
	// Open our jsonFile
	jsonFile, err := os.Open("data/pro-players.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully Opened pro-players.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal([]byte(byteValue), &proPlayers)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Successfully Loaded pro-players.json")

	for dotaId, twitchId := range dotaToTwitchMap {
		reverseLookup[twitchId] = dotaId
	}
}

var twitchClient = api.CreateTwitchClient()

func fetchData(ids []string) []helix.Stream {

	// Add stream data
	twitchResp, err := twitchClient.GetStreams(&helix.StreamsParams{
		GameIDs: []string{api.DotaGameId},
		UserIDs: ids,
	})
	if err != nil {
		log.Error(err)
	}

	if twitchResp.Error != "" {
		log.Error(err)
	}

	return twitchResp.Data.Streams
}

func LookupPlayers(list []dota.GameList) []helix.Stream {

	ids := dota.ExtractUserIds(list)

	twitchIds := lookupTwitchIds(ids)
	log.Info("Found twitch user ids ", twitchIds)

	if len(twitchIds) == 0 {
		return []helix.Stream{}
	}

	data := fetchData(twitchIds)

	return data
}

func lookupTwitchIds(userIds []int) []string {
	twitchUserIds := make([]string, 0)

	for _, id := range userIds {
		if val, ok := dotaToTwitchMap[id]; ok {
			twitchUserIds = append(twitchUserIds, val)
			//do something here
		}
	}

	return twitchUserIds
}

func AddStreamInfo(games *dota.TopLiveGames) *dota.TopLiveGames {

	steams := LookupPlayers(games.GameList)

	for i, game := range games.GameList {

		for playerKey, player := range game.Players {

			if twitchId, ok := dotaToTwitchMap[player.AccountID]; ok {

				//do something here
				for _, stream := range steams {
					if twitchId == stream.UserID {
						log.Info("Found stream for user ", player.AccountID)
						games.GameList[i].Players[playerKey].Stream = stream
					}
				}
			}

			if leaderboard, ok := api.DotaPlayers[player.AccountID]; ok {
				log.Info("Found DotaPlayers user id ", player.AccountID)
				games.GameList[i].Players[playerKey].LeaderboardRank = leaderboard.SteamAccount.SeasonLeaderboardRank
			}
		}
	}

	return games
}
