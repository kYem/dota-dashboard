package config

import "os"

func LoadConfig() ValveApiConfig {
	return ValveApiConfig{
		"api.steampowered.com",
		"https",
		os.Getenv("STEAM_API_KEY"),
		"IDOTA2Match_570",
	}
}

type ValveApiConfig struct {
	Hostname string
	Schema   string
	ApiKey   string
	DotaGame string
}
