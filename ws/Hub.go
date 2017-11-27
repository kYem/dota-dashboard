package ws

import (
	"github.com/satori/go.uuid"
	"log"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"github.com/kYem/dota-dashboard/dota"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)


// Store holds the collection of users connected through websocket
type Store struct {
	Users []*User
	Channels map[string]map[string]*User
	pubSubConn *redis.PubSubConn
	sync.Mutex
}

type User struct {
	ID   string
	conn *websocket.Conn
}

func (s *Store) newUser(conn *websocket.Conn) *User {
	u := &User{
		ID:   uuid.NewV4().String(),
		conn: conn,
	}
	log.Printf("Adding new user %s\n", u.ID)
	s.Lock()
	defer s.Unlock()

	s.Users = append(s.Users, u)
	return u
}

func (s *Store) Subscribe(u *User, channelName string) error {

	if _, ok := s.Channels[channelName]; ok {
		fmt.Sprintf("Channel %s already exists \n", channelName)
	} else {
		fmt.Sprintf("Creating empty channel %s \n", channelName)
		s.Channels[channelName] = make(map[string]*User)
	}

	if conErr := gPubSubConn.Subscribe(channelName); conErr != nil {
		return conErr
	}
	log.Println("Subscribed to live info for " + channelName)

	s.Lock()
	defer s.Unlock()
	s.Channels[channelName][u.ID] = u

	return nil
}


func (s *Store) findAndDeliver(channel string, content string) {

	var match dota.ApiLiveMatch
	if err := json.Unmarshal([]byte(content), &match); err != nil {
		log.Println(err)
	}

	wsResp := ApiMatchResponse{
		Event: channel,
		Data: match,
		Success: true,
	}

	if _, ok := s.Channels[channel]; ok {

		for _, u := range s.Channels[channel] {

			if err := u.conn.WriteJSON(wsResp); err != nil {
				log.Printf("error on message delivery through ws. e: %s\n", err)
				go s.removeUser(u)
			} else {
				log.Printf("user %s found at our store, message sent\n", u.ID)
			}
			return
		}
	} else {
		log.Printf("Channel %s not found at our store\n", channel)
	}

}


func (s *Store) removeUser(u *User) {

	s.Lock()
	defer s.Unlock()

	for i, storeUser := range s.Users {
		if storeUser == u {
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
			log.Printf("Removed user %s from store \n", u.ID)
		}
	}

	// Now remove from channels
	for channelName := range s.Channels {

		delete(s.Channels[channelName], u.ID)
		log.Printf("Remove user %s subscribtion %s\n", u.ID, channelName)

		remaining := len(s.Channels[channelName])
		log.Printf("channel %s have %d subs remaining\n", channelName, remaining)

		if remaining == 0 {
			delete(s.Channels, channelName)
			log.Printf("channel %s have been removed\n", channelName)
			if err := s.pubSubConn.Unsubscribe(channelName); err != nil {
				log.Printf("Failed to remove server subscription %s, err: %s\n", channelName, err)
			}
		}
	}
}
