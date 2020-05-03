package main

import (
	"encoding/json"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/cache"
	"github.com/kYem/dota-dashboard/dota"
	"github.com/kYem/dota-dashboard/storage"
	"github.com/kYem/dota-dashboard/stream"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var StratzClient = api.NewStratzClient("")

func SetDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("access-control-allow-credentials", "true")
	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func HomePage(w http.ResponseWriter, req *http.Request) {

	resp := api.SteamApi.GetTopLiveGames("1")

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(body))
}



func LiveGames(w http.ResponseWriter, req *http.Request) {

	SetDefaultHeaders(w)

	c := cache.Cache
	cacheKey := "live_games"
	cacheItem, err := c.Get(cacheKey)

	if err == nil {
		io.WriteString(w, cacheItem)
		return
	}


	partner := req.URL.Query().Get("partner")
	if partner == "" {
		partner = "0"
	}
	resp := api.SteamApi.GetTopLiveGames(partner)
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received api %d\n", resp.StatusCode)
		http.Error(w, "Steam api is down", 500)
		return
	}

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()

	var liveGames dota.TopLiveGames
	if err := json.NewDecoder(resp.Body).Decode(&liveGames); err != nil {
		log.Println(err)
	}

	stream.AddStreamInfo(&liveGames)
	for i, game := range liveGames.GameList {
		for playerKey, player := range game.Players {
			liveGames.GameList[i].Players[playerKey].Hero = storage.HeroById(player.HeroID)
		}
	}

	// Send back
	b, err := json.Marshal(liveGames)
	if err != nil {
		panic(err)
		return
	}

	bodyString := string(b)

	c.Set(cacheKey, bodyString, 10)


	io.WriteString(w, bodyString)
}

func LiveGamesStats(w http.ResponseWriter, req *http.Request) {
	SetDefaultHeaders(w)
	serverSteamId := req.URL.Query().Get("server_steam_id")
	resp := api.SteamApi.GetRealTimeStats(serverSteamId)

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()

	var match dota.LiveMatch
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		log.Println(err)
	}

	for i, game := range match.Teams {
		for playerKey, player := range game.Players {
			match.Teams[i].Players[playerKey].Hero = storage.HeroById(player.HeroId)
		}
	}
	// Send back
	json.NewEncoder(w).Encode(match)
}

func LeagueGames(w http.ResponseWriter, _ *http.Request) {
	SetDefaultHeaders(w)
	resp := api.SteamApi.GetLiveLeagueGames()

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()

	var leagueGames dota.LeagueGameResult
	if err := json.NewDecoder(resp.Body).Decode(&leagueGames); err != nil {
		log.Println(err)
	}
	err := json.NewEncoder(w).Encode(leagueGames)
	if err != nil {
		log.Println("Failed to Encode LeagueGames")
	}
}


func PassThrough(w http.ResponseWriter, req *http.Request) {
	SetDefaultHeaders(w)
	region := req.URL.Query().Get("region")
	skip := req.URL.Query().Get("start")

	i, err := strconv.ParseInt(skip, 10, 32)

	var result int32 = 0
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		result = int32(i)
	}

	resp := StratzClient.GetSeasonLeaderBoard(region, result, 100)

	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		http.Error(w, "Error contacting third party api", 500)
		return
	}
}