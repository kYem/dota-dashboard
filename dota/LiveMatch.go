package dota

import (
	"encoding/json"
	"strconv"
)

type SteamDotaPlayer struct {
	AccountId    int     `json:"accountid"`
	PlayerId     int     `json:"playerid"`
	HeroId       int     `json:"heroid"`
	PlayerStats
}

type PlayerStats struct {
	Name         string  `json:"name"`
	Team         int     `json:"team"`
	Level        int     `json:"level"`
	KillCount    int     `json:"kill_count"`
	DeathCount   int     `json:"death_count"`
	AssistsCount int     `json:"assists_count"`
	DeniesCount  int     `json:"denies_count"`
	LhCount      int     `json:"lh_count"`
	Gold         int     `json:"gold"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	NetWorth     int     `json:"net_worth"`
}

// Convert to DotaPlayer on encode, normalize json keys to underscore
func (u *SteamDotaPlayer) MarshalJSON() ([]byte, error) {
	var a = Player(*u)
	return json.Marshal(&a)
}

type Player struct {
	AccountId    int     `json:"account_id"`
	PlayerId     int     `json:"player_id"`
	HeroId       int     `json:"hero_id"`
	PlayerStats
}

type TeamDetails struct {
	TeamNumber int    `json:"team_number"`
	TeamID     int    `json:"team_id"`
	TeamName   string `json:"team_name"`
	TeamTag    string `json:"team_tag"`
	TeamLogo   int64    `json:"team_logo"`
	Score      int    `json:"score"`
	NetWorth   int    `json:"net_worth"`
}

type Team struct {
	TeamDetails
	Players    []SteamDotaPlayer `json:"players"`
}

type ApiTeam struct {
	TeamDetails
	Players    []Player `json:"players"`
}

func (team *ApiTeam) MarshalJSON() ([]byte, error) {
	type Alias ApiTeam
	return json.Marshal(&struct {
		TeamLogo string `json:"team_logo"`
		*Alias
	}{
		TeamLogo: strconv.FormatInt(team.TeamLogo, 10) ,
		Alias:   (*Alias)(team),
	})
}

type Building struct {
	Team      int  `json:"team"`
	Heading   float64  `json:"heading"`
	Type      int  `json:"type"`
	Lane      int  `json:"lane"`
	Tier      int  `json:"tier"`
	X         float64  `json:"x"`
	Y         float64  `json:"y"`
	Destroyed bool `json:"destroyed"`
}

type SteamDotaMatch struct {
	MatchId int64 `json:"matchid"`
	MatchDetails
}

func (match *SteamDotaMatch) MarshalJSON() ([]byte, error) {
	type Alias SteamDotaMatch
	return json.Marshal(&struct {
		MatchId string `json:"match_id"`
		ServerSteamID string `json:"server_steam_id"`
		*Alias
	}{
		MatchId: strconv.FormatInt(match.MatchId, 10) ,
		ServerSteamID: strconv.FormatInt(match.ServerSteamID, 10),
		Alias:   (*Alias)(match),
	})
}

type Match struct {
	MatchId string `json:"match_id"`
	MatchDetails
}

type MatchDetails struct {
	ServerSteamID int64 `json:"server_steam_id"`
	Timestamp     int   `json:"timestamp"`
	GameTime      int   `json:"game_time"`
	GameMode      int   `json:"game_mode"`
	LeagueID      int   `json:"league_id"`
	LeagueNodeID  int   `json:"league_node_id"`
	GameState     int   `json:"game_state"`
}

type GraphData struct {
	GraphGold []int `json:"graph_gold"`
}

type LiveMatch struct {
	Match *SteamDotaMatch `json:"match"`
	Teams []Team `json:"teams"`
	Buildings []Building `json:"buildings"`
	GraphData GraphData `json:"graph_data"`
	DeltaFrame bool `json:"delta_frame"`
}

type ApiLiveMatch struct {
	Match *Match `json:"match"`
	Teams []ApiTeam `json:"teams"`
	Buildings []Building `json:"buildings"`
	GraphData GraphData `json:"graph_data"`
	DeltaFrame bool `json:"delta_frame"`
}
