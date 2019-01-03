package dota

import "time"

type ProPlayer struct {
	AccountID       int         `json:"account_id"`
	SteamId         string      `json:"steamid"`
	Avatar          string      `json:"avatar"`
	AvatarMedium    string      `json:"avatarmedium"`
	AvatarFull      string      `json:"avatarfull"`
	ProfileUrl      string      `json:"profileurl"`
	PersonaName     string      `json:"personaname"`
	LastLogin       interface{} `json:"last_login"`
	FullHistoryTime time.Time   `json:"full_history_time"`
	Cheese          int         `json:"cheese"`
	FhUnavailable   bool        `json:"fh_unavailable"`
	LocCountryCode  string      `json:"loccountrycode"`
	LastMatchTime   time.Time   `json:"last_match_time"`
	Name            string      `json:"name"`
	CountryCode     string      `json:"country_code"`
	FantasyRole     int         `json:"fantasy_role"`
	TeamID          int         `json:"team_id"`
	TeamName        string      `json:"team_name"`
	TeamTag         string      `json:"team_tag"`
	IsLocked        bool        `json:"is_locked"`
	IsPro           bool        `json:"is_pro"`
	LockedUntil     int         `json:"locked_until"`
}