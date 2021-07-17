package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type StratzClient struct {
	hostname string
}

type StratzPlayer struct {
	SteamAccount struct {
		ID                          int    `json:"id"`
		SteamID                     int    `json:"steamId"`
		LastActiveTime              string `json:"lastActiveTime"`
		ProfileURI                  string `json:"profileUri"`
		TimeCreated                 int    `json:"timeCreated"`
		CityID                      int    `json:"cityId"`
		CommunityVisibleState       int    `json:"communityVisibleState"`
		Name                        string `json:"name"`
		Avatar                      string `json:"avatar"`
		PrimaryClanID               int64  `json:"primaryClanId"`
		SoloRank                    int    `json:"soloRank"`
		PartyRank                   int    `json:"partyRank"`
		IsDotaPlusSubscriber        bool   `json:"isDotaPlusSubscriber"`
		DotaPlusOriginalStartDate   int    `json:"dotaPlusOriginalStartDate"`
		IsAnonymous                 bool   `json:"isAnonymous"`
		IsStratzAnonymous           bool   `json:"isStratzAnonymous"`
		SeasonRank                  int    `json:"seasonRank"`
		SeasonLeaderboardRank       int    `json:"seasonLeaderboardRank"`
		SeasonLeaderboardDivisionID int    `json:"seasonLeaderboardDivisionId"`
		ProSteamAccount             struct {
			SteamID         int      `json:"steamId"`
			Name            string   `json:"name"`
			RealName        string   `json:"realName"`
			FantasyRole     int      `json:"fantasyRole"`
			TeamID          int      `json:"teamId"`
			Sponsor         string   `json:"sponsor"`
			IsLocked        bool     `json:"isLocked"`
			IsPro           bool     `json:"isPro"`
			TotalEarnings   int      `json:"totalEarnings"`
			Roles           int      `json:"roles"`
			Aliases         []string `json:"aliases"`
			Statuses        int      `json:"statuses"`
			TwitterLink     string   `json:"twitterLink"`
			TwitchLink      string   `json:"twitchLink"`
			InstagramLink   string   `json:"instagramLink"`
			FacebookLink    string   `json:"facebookLink"`
			SignatureHeroes []string `json:"signatureHeroes"`
			Countries       []string `json:"countries"`
			TiWins          int      `json:"tiWins"`
			IsTIWinner      bool     `json:"isTIWinner"`
		} `json:"proSteamAccount"`
		SmurfFlag         int `json:"smurfFlag"`
		SmurfCheckDate    int `json:"smurfCheckDate"`
		LastMatchDateTime int `json:"lastMatchDateTime"`
		LastMatchRegionID int `json:"lastMatchRegionId"`
	} `json:"steamAccount"`
	LastUpdateDateTime int `json:"lastUpdateDateTime"`
	RankVariance       int `json:"rankVariance"`
	MatchCount         int `json:"matchCount"`
	CoreCount          int `json:"coreCount"`
	SupportCount       int `json:"supportCount"`
	Imp                int `json:"imp"`
	Heroes             []struct {
		HeroID    int `json:"heroId"`
		WinCount  int `json:"winCount"`
		LossCount int `json:"lossCount"`
	} `json:"heroes"`
	Activity int `json:"activity"`
}

type LeaderBoardDivisionResponse struct {
	LeaderBoardDivisionID int            `json:"leaderBoardDivisionId"`
	Players               []StratzPlayer `json:"players"`
	PlayerCount           int            `json:"playerCount"`
}

var DotaPlayers = map[int]StratzPlayer{}

func init() {

	client := NewStratzClient("")

	urls := getLeaderboardUrls(client)

	go addUserDetails(urls)

}

func addUserDetails(urls []string) {
	results := BoundedParallelGet(urls, 2)

	for _, result := range results {
		data := processResponse(result.err, &result.res)
		for _, player := range data.Players {
			DotaPlayers[player.SteamAccount.ID] = player
		}
	}

	log.Printf("Added users count Stratz %d\n", len(DotaPlayers))
}

func getLeaderboardUrls(client *StratzClient) []string {
	var urls []string
	var take int32 = 100
	for i := 0; i <= 3; i++ {
		var start int32 = 0
		for ; start < 300; {
			url := client.GetSeasonLeaderBoardUrl(strconv.Itoa(i), start, take)
			urls = append(urls, url)
			start += take
		}
	}
	return urls
}

func NewStratzClient(hostname string) *StratzClient {
	if hostname == "" {
		hostname = "https://api.stratz.com/api/v1"
	}
	return &StratzClient{hostname: hostname}
}

func (client *StratzClient) GetSeasonLeaderBoardUrl(region string, skip int32, take int32) string {

	return fmt.Sprintf(
		"%s/Player/seasonLeaderBoard?&leaderBoardDivision=%s&skip=%d&take=%d",
		client.hostname,
		region,
		skip,
		take,
	)
}

func (client *StratzClient) GetSeasonLeaderBoard(region string, skip int32, take int32) LeaderBoardDivisionResponse {

	url := client.GetSeasonLeaderBoardUrl(region, skip, take)
	resp, err := http.Get(url)

	return processResponse(err, resp)
}

func processResponse(err error, resp *http.Response) LeaderBoardDivisionResponse {
	if err != nil {
		log.Printf("GetSeasonLeaderBoard: failed: %s", err)
	}
	var leaderBoardResp LeaderBoardDivisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&leaderBoardResp); err != nil {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to load Leader boards, %s, %s", err.Error(), string(body))
		}
		closeError := resp.Body.Close()
		if closeError != nil {
			log.Println("Error closing", closeError)
		}
	}
	return leaderBoardResp
}
