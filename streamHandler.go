package main

import (
	"encoding/json"
	"github.com/kYem/dota-dashboard/api"
	"github.com/nicklaw5/helix"
	"net/http"
)

var twitchClient = api.CreateTwitchClient()

func Streams(w http.ResponseWriter, req *http.Request) {
	SetDefaultHeaders(w)
	resp, err := twitchClient.GetStreams(&helix.StreamsParams{
		First:    10,
		Language: []string{"en"},
		GameIDs: []string{api.DotaGameId},
	})
	if err != nil {
		// handle error
		http.Error(w, err.Error(), 400)
		return
	}
	err = json.NewEncoder(w).Encode(resp.Data.Streams)
	if err != nil {
		// handle error
		http.Error(w, err.Error(), 400)
		return
	}
}


