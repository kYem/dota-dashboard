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
	"strings"
	"errors"
	"time"
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
	Channels []string
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

func (s *Store) newUser(conn *websocket.Conn) *User {
	u := &User{
		ID:   uuid.NewV4().String(),
		conn: conn,
		Channels: make([]string, 0, 1),
	}
	log.Printf("Adding new user %s\n", u.ID)
	s.Lock()
	defer s.Unlock()

	s.Users = append(s.Users, u)
	return u
}

func (s *Store) Subscribe(u *User, channelName string) error {

	if _, ok := s.Channels[channelName]; !ok {
		fmt.Sprintf("Creating empty channel %s \n", channelName)
		s.Channels[channelName] = make(map[string]*User)

		if conErr := gPubSubConn.Subscribe(channelName); conErr != nil {
			return conErr
		}
		log.Printf("subscribed to live info for %s", channelName)
	}

	s.Channels[channelName][u.ID] = u
	u.addChannel(channelName)

	return nil
}

func (s *Store) Unsubscribe(u *User, channelName string) {

	if _, ok := s.Channels[channelName]; !ok {
		log.Printf("Trying to unsubscribe from non existant channel %s \n", channelName)
	}

	delete(s.Channels[channelName], u.ID)
	u.removeChannel(channelName)
	log.Printf("store remove user %s subscribtion %s\n", u.ID, channelName)

	remaining := len(s.Channels[channelName])
	log.Printf("channel %s have %d subs remaining\n", channelName, remaining)

	if remaining == 0 {
		s.removeChannel(channelName)
	}
}

// Only allow single sub on live match
func (s *Store) SubscribeMatch(u *User, channelName string) error {

	// Make sure we are dealing with live match sub
	if ok := isLiveMatchChannel(channelName); !ok {
		return errors.New("channel name must start with " + channelLiveMatchPrefix + " received " + channelName)
	}

	// Unsubscribe from other channel
	var isMatchChannel bool
	for _, name := range u.Channels {

		// Already subscribed to this channel
		if name == channelName {
			continue
		}

		if isMatchChannel = isLiveMatchChannel(name); isMatchChannel {
			s.Unsubscribe(u, name)
			u.removeChannel(channelName)
		}
	}

	return s.Subscribe(u, channelName)
}

func isLiveMatchChannel(channelName string) bool {
	return strings.HasPrefix(channelName, channelLiveMatchPrefix)
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

		log.Printf("Broadcasting to %s, user %d \n", channel, len(s.Channels[channel]))
		start := time.Now()
		for _, u := range s.Channels[channel] {

			if err := u.conn.WriteJSON(wsResp); err != nil {
				log.Printf("error on message delivery through ws. e: %s\n", err)
				go s.removeUser(u)
			} else {
				log.Printf("user %s found at our store, message sent\n", u.ID)
			}
		}
		elapsed := time.Since(start)
		log.Printf("Delivered in took %s", elapsed)
	} else {
		log.Printf("Channel %s not found at our store\n", channel)
	}
}

func (s *Store) removeUser(u *User) {

	s.Lock()
	defer s.Unlock()

	s.removeUserFromChannels(u)

	for i, storeUser := range s.Users {
		if storeUser == u {
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
			log.Printf("Removed user %s from store \n", u.ID)
		}
	}
}

func (s *Store) removeUserFromChannels(u *User) {
	// Now remove from channels
	for _, channelName := range u.Channels {
		s.Unsubscribe(u, channelName)
	}

	u.Channels = make([]string, 0, 1)
}

func (s *Store) removeChannel(channelName string) {
	delete(s.Channels, channelName)
	log.Printf("channel %s have been removed\n", channelName)
	if err := s.pubSubConn.Unsubscribe(channelName); err != nil {
		log.Printf("Failed to remove server subscription %s, err: %s\n", channelName, err)
	}
}
