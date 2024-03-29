package main

import (
	"encoding/json"
	"errors"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/cache"
	"github.com/kYem/dota-dashboard/dota"
	"github.com/kYem/dota-dashboard/storage"
	"github.com/kYem/dota-dashboard/stream"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func JSONError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	resp := ErrorResponse{
		Error: err.Error(),
		Code:  code,
	}

	encodeErr := json.NewEncoder(w).Encode(resp)
	if encodeErr != nil {
		log.Println(encodeErr, err, code)
		return
	}
}

var DefaultError = errors.New("there was a problem completing your request")

func HomePage(w http.ResponseWriter, _ *http.Request) {

	resp := api.SteamApi.GetTopLiveGames("1")

	if resp.Body == nil {
		JSONError(w, DefaultError, 400)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_, _ = io.WriteString(w, string(body))
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
	if resp.StatusCode != http.StatusOK || resp.Body == nil {
		log.Printf("Received api %d\n", resp.StatusCode)
		JSONError(w, errors.New("steam API is down"), 500)
		return
	}

	defer resp.Body.Close()

	var liveGames dota.TopLiveGames
	if err := json.NewDecoder(resp.Body).Decode(&liveGames); err != nil {
		log.Println(err)
	}

	stream.AddPlayerInfo(&liveGames)
	for i, game := range liveGames.GameList {
		for playerKey, player := range game.Players {
			liveGames.GameList[i].Players[playerKey].Hero = storage.HeroById(player.HeroID)
		}
	}

	// Send back
	b, err := json.Marshal(liveGames)
	if err != nil {
		log.Error(err)
		JSONError(w, DefaultError, 500)
		return
	}

	bodyString := string(b)

	err = c.Set(cacheKey, bodyString, 10)
	if err != nil {
		log.Error(err)
	}

	_, _ = io.WriteString(w, bodyString)
}

func LiveGamesStats(w http.ResponseWriter, req *http.Request) {
	SetDefaultHeaders(w)
	serverSteamId := req.URL.Query().Get("server_steam_id")
	match, err := api.SteamApi.GetRealTimeStats(serverSteamId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	for i, game := range match.Teams {
		for playerKey, player := range game.Players {
			match.Teams[i].Players[playerKey].Hero = storage.HeroById(player.HeroId)
		}
	}
	// Send back
	err = json.NewEncoder(w).Encode(match)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
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
