package config

func LoadConfig() ValveApiConfig {
	return ValveApiConfig{
		"api.steampowered.com",
		"https",
		"673EF461E2406D27FDE71A7E20DFBAF1",
		"IDOTA2Match_570",
	}
}

type ValveApiConfig struct {
	Hostname string
	Schema   string
	ApiKey   string
	DotaGame string
}