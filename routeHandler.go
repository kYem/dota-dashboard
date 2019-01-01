package main

import (
	"net/http"
	"io/ioutil"
	"io"
	"github.com/kYem/dota-dashboard/dota"
	"encoding/json"
	"log"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/config"
	"github.com/kYem/dota-dashboard/cache"
)

func SetDefaultHeaders(w http.ResponseWriter) {
	w.Header().Add("access-control-allow-credentials", "true")
	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func HomePage(w http.ResponseWriter, req *http.Request) {

	client := api.GetClient(config.LoadConfig())

	resp := client.GetTopLiveGames("1")

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
	client := api.GetClient(config.LoadConfig())
	resp := client.GetTopLiveGames(partner)
	if resp.StatusCode != http.StatusOK {
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
	client := api.GetClient(config.LoadConfig())
	resp := client.GetRealTimeStats(serverSteamId)

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()

	var match dota.LiveMatch
	if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
		log.Println(err)
	}
	// Send back
	json.NewEncoder(w).Encode(match)
}

func LeagueGames(w http.ResponseWriter, req *http.Request) {
	SetDefaultHeaders(w)
	client := api.GetClient(config.LoadConfig())
	resp := client.GetLiveLeagueGames()

	if resp.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	defer resp.Body.Close()

	var leagueGames dota.LeagueGameResult
	if err := json.NewDecoder(resp.Body).Decode(&leagueGames); err != nil {
		log.Println(err)
	}
	json.NewEncoder(w).Encode(leagueGames)
}

