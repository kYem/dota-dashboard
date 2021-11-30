package main

import (
	"encoding/json"
	"github.com/kYem/dota-dashboard/api"
	"net/http"
)

func Streams(w http.ResponseWriter, req *http.Request) {
	SetDefaultHeaders(w)
	resp, err := api.TwitchClient.GetStreams([]string{}, 10)
	if err != nil {
		JSONError(w, err, 500)
		return
	}
	err = json.NewEncoder(w).Encode(resp.Data.Streams)
	if err != nil {
		// handle error
		JSONError(w, err, 500)
		return
	}
}
