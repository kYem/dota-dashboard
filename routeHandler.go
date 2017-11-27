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
)

func SetDefaultHeaders(w http.ResponseWriter) {
	w.Header().Add("access-control-allow-credentials", "true")
	w.Header().Add("access-control-allow-origin", "http://dotatv.com:3000")
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
	partner := req.URL.Query().Get("partner")
	if partner == "" {
		partner = "0"
	}
	client := api.GetClient(config.LoadConfig())
	resp := client.GetTopLiveGames(partner)

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

