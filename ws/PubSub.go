package ws

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/kYem/dota-dashboard/dota"
	"github.com/kYem/dota-dashboard/cache"
)

const (
	channelLiveMatchPrefix = "dota_live_match."
)

var (
	gStore      *Store
	gPubSubConn *redis.PubSubConn
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

type ApiMatchResponse struct {
	Event string `json:"event"`
	Data dota.ApiLiveMatch `json:"data"`
	Success bool `json:"success"`
}

func Init() {

	gPubSubConn = GetConn()
	gStore = &Store{
		Users: make([]*User, 0, 1),
		Channels: make(map[string]map[string]*User),
		pubSubConn: gPubSubConn,
	}

	go DeliverMessages()
}

var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrader error %s\n" + err.Error())
		return
	}
	u := gStore.newUser(conn)

	go u.writePump()

	for {

		// receive JSON type T
		var data WsRequest
		if err = conn.ReadJSON(&data); err != nil {
			log.Println("Can't receive", err.Error())
			gStore.removeUser(u)
			break
		}

		err := gStore.SubscribeMatch(u, channelLiveMatchPrefix+data.Params.ServerSteamID)
		if err != nil {
			log.Println(err)
		}

	}
}

func GetConn() *redis.PubSubConn  {

	c := cache.Pool.Get()

	log.Printf("Active connections %d", cache.Pool.ActiveCount())
	pubSub := &redis.PubSubConn{Conn: c}

	if pubSub == nil {
		fmt.Println("Failed to get pubSub con")
	}

	return pubSub
}


func DeliverMessages() {

	for {
		switch v := gPubSubConn.Receive().(type) {
		case redis.Message:
			gStore.findAndDeliver(v.Channel, string(v.Data))
		case redis.Subscription:
			log.Printf("subscription message: %s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Println("error pub/sub on connection, delivery has stopped")
			log.Printf("error %s", v.Error())
			return
		}
	}
}
