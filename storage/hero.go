package storage

import (
	"encoding/json"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/dota"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const (
	// cdn            = "https://steamcdn-a.akamaihd.net/apps/dota2/images/heroes/"
	cdn            = "https://cdn.cloudflare.steamstatic.com/apps/dota2/images/dota_react/heroes/"
	imageSize      = ".png"
	heroesFilename = "data/heroes.json"
)

var heroMap = map[int]dota.Hero {
	0: {
		Id: 0,
		Name: "unknown",
		Image: "/images/heroes/unknown-hero.jpeg",
	},
}


func init() {
	populateHeroMap()
}

func loadHeroes() {
	heroes, err := getHeroes()
	if err != nil {
		return
	}
	for _, hero := range heroes {
		heroMap[hero.Id] = dota.Hero{
			Id:    hero.Id,
			Name:  hero.Name,
			Image: cdn + strings.Replace(hero.Name, "npc_dota_hero_", "", 1) + imageSize,
		}
	}
}

func getHeroes() ([]dota.HeroBasic, error) {
	heroes, err := api.SteamApi.GetHeroes()
	if err == nil {
		writeHeroJson(heroes)
		return heroes, nil
	}

	return loadHeroFile()
}

func writeHeroJson(heroes []dota.HeroBasic) {
	content, err := json.MarshalIndent(heroes, "", " ")
	if err == nil {
		_ = ioutil.WriteFile(heroesFilename, content, 0644)
	}
}


func populateHeroMap() {
	loadHeroes()
}


func loadHeroFile() ([]dota.HeroBasic, error) {

	var heroes []dota.HeroBasic

	// Open our jsonFile
	jsonFile, err := os.Open(heroesFilename)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	log.Info("Successfully Opened '%s'", heroesFilename)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &heroes)
	if err != nil {
		return nil, err
	}
	log.Info("Successfully Loaded '%s'", heroesFilename)
	return heroes, nil
}

func HeroById(id int) dota.Hero {
	return heroMap[id]
}
