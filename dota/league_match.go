package dota


type Game struct {
	Players []struct {
		AccountID int    `json:"account_id"`
		Name      string `json:"name"`
		HeroID    int    `json:"hero_id"`
		Team      int    `json:"team"`
	} `json:"players"`
	DireTeam struct {
		TeamName string `json:"team_name"`
		TeamID   int    `json:"team_id"`
		TeamLogo int    `json:"team_logo"`
		Complete bool   `json:"complete"`
	} `json:"dire_team"`
	LobbyID           int64  `json:"lobby_id"`
	MatchID           int64  `json:"match_id"`
	Spectators        int    `json:"spectators"`
	SeriesID          int    `json:"series_id"`
	GameNumber        int    `json:"game_number"`
	LeagueID          int    `json:"league_id"`
	StreamDelayS      int    `json:"stream_delay_s"`
	RadiantSeriesWins int    `json:"radiant_series_wins"`
	DireSeriesWins    int    `json:"dire_series_wins"`
	SeriesType        int    `json:"series_type"`
	LeagueSeriesID    int    `json:"league_series_id"`
	LeagueGameID      int    `json:"league_game_id"`
	StageName         string `json:"stage_name"`
	LeagueTier        int    `json:"league_tier"`
	Scoreboard        struct {
		Duration           float64 `json:"duration"`
		RoshanRespawnTimer int     `json:"roshan_respawn_timer"`
		Radiant            struct {
			LeagueTeam
		} `json:"radiant"`
		Dire struct {
			LeagueTeam
		} `json:"dire"`
	} `json:"scoreboard"`
}

type LeagueTeam struct {
	Score         int `json:"score"`
	TowerState    int `json:"tower_state"`
	BarracksState int `json:"barracks_state"`
	Picks         []struct {
		HeroID int `json:"hero_id"`
	} `json:"picks"`
	Bans []struct {
		HeroID int `json:"hero_id"`
	} `json:"bans"`
	Players []struct {
		PlayerSlot       int     `json:"player_slot"`
		AccountID        int     `json:"account_id"`
		HeroID           int     `json:"hero_id"`
		Kills            int     `json:"kills"`
		Death            int     `json:"death"`
		Assists          int     `json:"assists"`
		LastHits         int     `json:"last_hits"`
		Denies           int     `json:"denies"`
		Gold             int     `json:"gold"`
		Level            int     `json:"level"`
		GoldPerMin       int     `json:"gold_per_min"`
		XpPerMin         int     `json:"xp_per_min"`
		UltimateState    int     `json:"ultimate_state"`
		UltimateCooldown int     `json:"ultimate_cooldown"`
		Item0            int     `json:"item0"`
		Item1            int     `json:"item1"`
		Item2            int     `json:"item2"`
		Item3            int     `json:"item3"`
		Item4            int     `json:"item4"`
		Item5            int     `json:"item5"`
		RespawnTimer     int     `json:"respawn_timer"`
		PositionX        float64 `json:"position_x"`
		PositionY        float64 `json:"position_y"`
		NetWorth         int     `json:"net_worth"`
	} `json:"players"`
	Abilities []struct {
		AbilityID    int `json:"ability_id"`
		AbilityLevel int `json:"ability_level"`
	} `json:"abilities"`
}

type LeagueGames struct {
	Games []Game `json:"games"`
}

type LeagueGameResult struct {
	Result struct {
		Games []Game `json:"games"`
	} `json:"result"`
}