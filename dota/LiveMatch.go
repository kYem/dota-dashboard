package dota


type SteamDotaPlayer struct {
	Accountid    int     `json:"accountid"`
	Playerid     int     `json:"playerid"`
	Heroid       int     `json:"heroid"`
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
}

type DotaPlayer struct {
	AccountId    int     `json:"account_id"`
	PlayerId     int     `json:"player_id"`
	HeroId       int     `json:"hero_id"`
	PlayerStats
}

type Team struct {
	TeamNumber int    `json:"team_number"`
	TeamID     int    `json:"team_id"`
	TeamName   string `json:"team_name"`
	TeamLogo   string `json:"team_logo"`
	Score      int    `json:"score"`
	Players    []SteamDotaPlayer `json:"players"`
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
	ServerSteamID string `json:"server_steam_id"`
	Matchid       string `json:"matchid"`
	Timestamp     int    `json:"timestamp"`
	GameTime      int    `json:"game_time"`
	GameMode      int    `json:"game_mode"`
	LeagueID      int    `json:"league_id"`
}

type GraphData struct {
	GraphGold []int `json:"graph_gold"`
}

type LiveMatch struct {
	Match SteamDotaMatch `json:"match"`
	Teams []Team `json:"teams"`
	Buildings []Building `json:"buildings"`
	GraphData GraphData `json:"graph_data"`
	DeltaFrame bool `json:"delta_frame"`
}