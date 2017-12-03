package ws

import (
	//"time"
	"log"
	"github.com/gorilla/websocket"
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
		u.conn.Close()
	}()
	for {
		select {
		case wsResp := <-u.send:
			//u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.conn.WriteJSON(wsResp); err != nil {
				log.Printf("error on message delivery through ws. e: %s\n", err)
				gStore.removeUser(u)
			}
		}
	}
}

func (u *User) addChannel(channelName string) {
	u.Channels = append(u.Channels, channelName)
}

func (u *User) removeChannel(channelName string) {
	for i, channel := range u.Channels {
		if channel == channelName {
			u.Channels = append(u.Channels[:i], u.Channels[i+1:]...)
			log.Printf("removed user %s channel %s \n", u.ID, channel)
			break
		}
	}
}
