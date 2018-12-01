package dota

import (
	"encoding/json"
	"strconv"
)

type GameList struct {
	ActivateTime    int    `json:"activate_time"`
	DeactivateTime  int    `json:"deactivate_time"`
	ServerSteamID   int64  `json:"server_steam_id"`
	LobbyID         int64  `json:"lobby_id"`
	LeagueID        int    `json:"league_id"`
	LobbyType       int    `json:"lobby_type"`
	GameTime        int    `json:"game_time"`
	Delay           int    `json:"delay"`
	Spectators      int    `json:"spectators"`
	GameMode        int    `json:"game_mode"`
	AverageMmr      int    `json:"average_mmr"`
	MatchID         int64  `json:"match_id"`
	SeriesID        int    `json:"series_id"`
	TeamNameRadiant string `json:"team_name_radiant"`
	TeamNameDire    string `json:"team_name_dire"`
	TeamLogoRadiant int64  `json:"team_logo_radiant"`
	TeamLogoDire    int64  `json:"team_logo_dire"`
	TeamIDRadiant   int    `json:"team_id_radiant"`
	TeamIDDire      int    `json:"team_id_dire"`
	SortScore       int    `json:"sort_score"`
	LastUpdateTime  int    `json:"last_update_time"`
	RadiantLead     int    `json:"radiant_lead"`
	RadiantScore    int    `json:"radiant_score"`
	DireScore       int    `json:"dire_score"`
	Players         []struct {
		AccountID int `json:"account_id"`
		HeroID    int `json:"hero_id"`
	} `json:"players"`
	BuildingState int `json:"building_state"`
}
type TopLiveGames struct {
	GameList []GameList `json:"game_list"`
}

func (match *GameList) MarshalJSON() ([]byte, error) {
	type Alias GameList
	return json.Marshal(&struct {
		ServerSteamID string `json:"server_steam_id"`
		LobbyID string `json:"lobby_id"`
		MatchID string `json:"match_id"`
		TeamLogoRadiant string  `json:"team_logo_radiant"`
		TeamLogoDire    string  `json:"team_logo_dire"`

		*Alias
	}{
		ServerSteamID: strconv.FormatInt(match.ServerSteamID, 10),
		LobbyID: strconv.FormatInt(match.LobbyID, 10),
		MatchID: strconv.FormatInt(match.MatchID, 10),
		TeamLogoRadiant: strconv.FormatInt(match.TeamLogoRadiant, 10),
		TeamLogoDire: strconv.FormatInt(match.TeamLogoDire, 10),
		Alias:   (*Alias)(match),
	})
}