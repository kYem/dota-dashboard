package dota

type HeroBasic struct {
	Id       int     `json:"id"`
	Name 	string `json:"name"`
}

type HeroResult struct {
	Heroes []HeroBasic `json:"heroes"`
}

type GetHeroes struct {
	Result HeroResult `json:"result"`
}

type Hero struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}