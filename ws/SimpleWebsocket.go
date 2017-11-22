package ws

import (
	"github.com/kYem/dota-dashboard/dota"
	"fmt"
	"log"
	"encoding/json"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/config"
	"golang.org/x/net/websocket"
)

type LiveMatchParams struct {
	ServerSteamID string `json:"server_steam_id"`
}
type WsRequest struct {
	Event string
	Reference string
	Params LiveMatchParams
}

type WsError struct {
	Event string `json:"event"`
	Success bool `json:"success"`
	Error string
}

type MatchResponse struct {
	Event string `json:"event"`
	Data dota.LiveMatch `json:"data"`
	Success bool `json:"success"`
}

func Echo(ws *websocket.Conn) {
	defer ws.Close()
	fmt.Println("Client Connected")
	var err error

	for {

		// receive JSON type T
		var data WsRequest
		if err = websocket.JSON.Receive(ws, &data); err != nil {
			fmt.Println("Can't receive", err.Error())
			break
		}

		client := api.GetClient(config.LoadConfig())
		log.Println("Fetching live server info for " + data.Params.ServerSteamID)
		resp := client.GetRealTimeStats(data.Params.ServerSteamID)

		if resp.Body == nil {
			apiError := WsError{
				Event: data.Event + "." + data.Reference,
				Error: "Failed to get live match data from steam",
				Success: false,
			}
			if err = websocket.JSON.Send(ws, apiError); err != nil {
				break
			}
		}

		var match dota.LiveMatch
		if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
			log.Println(err)
		}

		resp.Body.Close()

		wsResp := MatchResponse{
			Event: data.Event + "." + data.Reference,
			Data: match,
			Success: true,
		}
		if err = websocket.JSON.Send(ws, wsResp); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

