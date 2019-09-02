package stream

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/dota"
	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

const (
	gorgc           = 56939869
	cancel          = 141690233
	Kingteka        = 242697694
	SexyBamboe      = 20321748
	PurgeGamers     = 66296404
	absolutHabibi   = 140873889
	JustCooman      = 175463659
	SyndereN        = 4281729
	KheZu           = 169025618
	Ghostik         = 82297531
	Mason           = 315657960
	Funkefal        = 101117815
	GunnarDotA2     = 126238768
	bobruhatv       = 86953944
	Mickee          = 152962063
	MickeeTwo       = 106755427
	Febby           = 112377459
	Mage            = 178366364
	SingSing        = 19757254
	Wagamama        = 32995405
	inboss1k        = 842068996
	eskobartv       = 246953032
	meepoha3ap      = 183602223
	eternalenvyy    = 43276219
	monkeysForever  = 86811043
	siractionslacks = 68186278
	roccodota       = 106932684
	universe        = 87276347
	noctisak47      = 101239422
	yoyuou          = 170773146
	lukiluki        = 117311875
	bububu          = 106381989
	threethree      = 86698277
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
	SingSing:      "21390470",

	// ybicanoooobov
	64607133:  "68950614",
	190769379: "68950614",
	199194550: "68950614",
	206967271: "68950614",

	Wagamama:        "24811779",
	inboss1k:        "431644708",
	eskobartv:       "140883424",
	meepoha3ap:      "85674261",
	eternalenvyy:    "26954716",
	monkeysForever:  "34932688",
	siractionslacks: "21379187",
	roccodota:       "65421010",
	universe:        "32556389",
	noctisak47:      "141414675",
	yoyuou:          "41727944",
	lukiluki:        "36945314",
	bububu:          "22573825",
	threethree:      "132521253",
}


var proPlayers []dota.ProPlayer

type twitchMap struct {
	SteamId     string `json:"steamId"`
	TwitchLogin string `json:"twitchLogin"`
	TwitchId    string `json:"twitchId"`
}

func init() {
	loadProPlayers()

	twitchPlayers := loadExtraTwitchPlayers()
	for _, twitchData := range twitchPlayers {
		i, err := strconv.Atoi(twitchData.SteamId)
		if err != nil {
			log.Fatal(err)
		}
		dotaToTwitchMap[i] = twitchData.TwitchId
	}

	log.Info("Loaded twitch player map: ", len(dotaToTwitchMap))
}

func loadProPlayers() {
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
	err = json.Unmarshal(byteValue, &proPlayers)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully Loaded pro-players.json")
}

func loadExtraTwitchPlayers() []twitchMap {
	// Open our jsonFile
	jsonFile, err := os.Open("data/twitch.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully Opened twitch.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	var twitchPlayers []twitchMap
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &twitchPlayers)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully Loaded twitch.json")

	return twitchPlayers
}

var twitchClient = api.CreateTwitchClient()

func fetchData(ids []string) []helix.Stream {

	// Add stream data
	twitchResp, err := twitchClient.GetStreams(&helix.StreamsParams{
		GameIDs: []string{api.DotaGameId},
		UserIDs: ids,
	})
	if err != nil {
		log.Error("ERROR from Twitch", err)
	}

	if twitchResp.Error != "" {
		log.Error(twitchResp.ErrorMessage)
	}

	return twitchResp.Data.Streams
}

func LookupPlayers(list []dota.GameList) []helix.Stream {

	ids := dota.ExtractUserIds(list)

	twitchIds := lookupTwitchIds(ids)

	if len(twitchIds) == 0 {
		return []helix.Stream{}
	}

	log.Info("Found twitch user ids", twitchIds, len(twitchIds))
	data := fetchData(twitchIds)

	return data
}

func lookupTwitchIds(userIds []int) []string {
	var twitchUserIds []string
	for _, id := range userIds {
		if val, ok := dotaToTwitchMap[id]; ok {
			twitchUserIds = append(twitchUserIds, val)
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
				games.GameList[i].Players[playerKey].LeaderboardRank = leaderboard.SteamAccount.SeasonLeaderboardRank
			}
		}
	}

	return games
}
