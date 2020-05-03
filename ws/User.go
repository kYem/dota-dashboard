package ws

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type User struct {
	ID   string
	conn *websocket.Conn
	Channels []string
	// Buffered channel of outbound messages.
	send chan *ApiMatchResponse
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (u *User) writePump() {
	defer func() {
		err := u.conn.Close()
		if err != nil {
			log.Error(err, "Closing Connection in WritePump FAILED")
		}
	}()
	for {
		select {
		case wsResp, ok := <-u.send:
			err := u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Errorf("SetWriteDeadline %v", err)
			}

			if !ok {
				// The hub closed the channel.
				err := u.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Error(err, "writePump closed hub")
				}
				return
			}
			if err := u.conn.WriteJSON(wsResp); err != nil {
				log.Errorf("error on message delivery through ws. e: %s\n", err)
				gStore.removeUser(u)
			}
		}
	}
}

func (u *User) readPump() {

	u.conn.SetReadLimit(maxMessageSize)
	_ = u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(
		func(string) error {
			_ = u.conn.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})

	for {
		// receive JSON type T
		var request Request
		if err := u.conn.ReadJSON(&request); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseNoStatusReceived) {
				log.Errorf("Unexpected close error: %v", err)
			}
			gStore.removeUser(u)
			break
		}

		err := gStore.SubscribeMatch(u, channelLiveMatchPrefix+request.Params.ServerSteamID)
		if err != nil {
			log.Errorf("Failed to subscribeMatch %v", err)
		}
	}

	err := u.conn.Close()
	if err != nil {
		log.Errorf("Failed to close user connection %v", err)
	}
}

func (u *User) addChannel(channelName string) {
	u.Channels = append(u.Channels, channelName)
}

func (u *User) removeChannel(channelName string) {
	for i, channel := range u.Channels {
		if channel == channelName {
			u.Channels = append(u.Channels[:i], u.Channels[i+1:]...)
			break
		}
	}
}
