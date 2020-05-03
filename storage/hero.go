package storage

import (
	"encoding/json"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/config"
	"github.com/kYem/dota-dashboard/dota"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

const cdn = "https://api.opendota.com/apps/dota2/images/heroes/"
const imageSize = "_sb.png"


var heroMap = map[int]dota.Hero {}

var steamApi = api.GetClient(config.LoadConfig())

func init() {
	populateHeroMap()
}

func loadHeroes() {
	heroes, err := steamApi.GetHeroes()
	if err != nil {
		heroes = loadHeroFile()
	}
	for _, hero := range heroes {
		heroMap[hero.Id] = dota.Hero{
			Id:    hero.Id,
			Name:  hero.Name,
			Image: cdn + strings.Replace(hero.Name, "npc_dota_hero_", "", 1) + imageSize,
		}
	}
}


func populateHeroMap() {
	loadHeroes()
}


var name = "data/heroes.json"
func loadHeroFile() []dota.HeroBasic {

	var heroes []dota.HeroBasic

	// Open our jsonFile
	jsonFile, err := os.Open(name)
	if err != nil {
		log.Errorf(err.Error())
		return heroes
	}

	log.Info("Successfully Opened '%s'", name)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &heroes)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Successfully Loaded '%s'", name)
	return heroes
}

func HeroById(id int) dota.Hero {
	return heroMap[id]
}
