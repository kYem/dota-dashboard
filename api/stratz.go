package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type StratzClient struct {
	hostname      string
}

type LeaderBoardDivisionResponse struct {
	LeaderBoardDivisionID int `json:"leaderBoardDivisionId"`
	Players               []struct {
		SteamAccount struct {
			ID                          int    `json:"id"`
			Name                        string `json:"name"`
			Avatar                      string `json:"avatar"`
			IsAnonymous                 bool   `json:"isAnonymous"`
			IsStratzAnonymous           bool   `json:"isStratzAnonymous"`
			SeasonRank                  int    `json:"seasonRank"`
			SeasonLeaderboardRank       int    `json:"seasonLeaderboardRank"`
			SeasonLeaderboardDivisionID int    `json:"seasonLeaderboardDivisionId"`
			ProSteamAccount             struct {
				SteamID       int    `json:"steamId"`
				Name          string `json:"name"`
				RealName      string `json:"realName"`
				FantasyRole   int    `json:"fantasyRole"`
				TeamID        int    `json:"teamId"`
				Sponsor       string `json:"sponsor"`
				IsLocked      bool   `json:"isLocked"`
				IsPro         bool   `json:"isPro"`
				TotalEarnings int    `json:"totalEarnings"`
			} `json:"proSteamAccount"`
		} `json:"steamAccount"`
		LastUpdateDateTime int `json:"lastUpdateDateTime"`
		RankVariance       int `json:"rankVariance"`
		Imp                int `json:"imp"`
		Activity           int `json:"activity"`
		MatchCount         int `json:"matchCount,omitempty"`
		CoreCount          int `json:"coreCount,omitempty"`
		SupportCount       int `json:"supportCount,omitempty"`
		Heroes             []struct {
			HeroID    int `json:"heroId"`
			WinCount  int `json:"winCount"`
			LossCount int `json:"lossCount"`
		} `json:"heroes,omitempty"`
	} `json:"players"`
	PlayerCount int `json:"playerCount"`
}


func NewStratzClient(hostname string) *StratzClient {
	if hostname == "" {
		hostname = "https://api.stratz.com/api/v1"
	}
	return &StratzClient{hostname: hostname}
}

func (client *StratzClient) GetSeasonLeaderBoard(region string, skip int32, take int32) LeaderBoardDivisionResponse {

	url := fmt.Sprintf(
		"%s/Player/seasonLeaderBoard?&leaderBoardDivision=%s",
		client.hostname,
		region,
	)
	resp, err := http.Get(url)

	if err != nil {
		log.Printf("GetTopLiveGames: failed: %s", err)
	}

	var leaderBoardResp LeaderBoardDivisionResponse
	if err := json.NewDecoder(resp.Body).Decode(&leaderBoardResp); err != nil {
		log.Println(err)
	}

	return leaderBoardResp
}
