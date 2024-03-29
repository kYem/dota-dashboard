package ws

import (
	"encoding/json"
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/kYem/dota-dashboard/api"
	"github.com/kYem/dota-dashboard/dota"
	"github.com/kYem/dota-dashboard/storage"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

// Store holds the collection of users connected through websocket
type Store struct {
	Users      []*User
	Channels   map[string]map[string]*User
	pubSubConn *redis.PubSubConn
	sync.Mutex
}

func (s *Store) newUser(conn *websocket.Conn) *User {
	u := &User{
		ID:       uuid.NewV4().String(),
		conn:     conn,
		Channels: make([]string, 0, 1),
		send:     make(chan *ApiMatchResponse),
	}
	s.Lock()
	defer s.Unlock()

	s.Users = append(s.Users, u)
	return u
}

func (s *Store) Subscribe(u *User, channelName string) error {

	s.Lock()
	if _, ok := s.Channels[channelName]; !ok {
		log.Infof("Creating empty channel %s \n", channelName)
		s.Channels[channelName] = make(map[string]*User)

		if conErr := gPubSubConn.Subscribe(channelName); conErr != nil {
			return conErr
		}
	}

	s.Channels[channelName][u.ID] = u
	u.addChannel(channelName)
	s.Unlock()

	return nil
}

func (s *Store) Unsubscribe(u *User, channelName string) {

	if _, ok := s.Channels[channelName]; !ok {
		log.Printf("Trying to unsubscribe from non existing channel %s \n", channelName)
	}

	delete(s.Channels[channelName], u.ID)
	u.removeChannel(channelName)
	log.Debugf("store remove user %s subscription %s\n", u.ID, channelName)

	remaining := len(s.Channels[channelName])
	log.Debugf("channel %s have %d subs remaining\n", channelName, remaining)

	if remaining == 0 {
		s.removeChannel(channelName)
	}
}

// SubscribeMatch Only allow single sub on live match
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

	err := s.sendInitialMatchData(u, channelName)
	if err != nil {
		log.Error(err)
	}

	return s.Subscribe(u, channelName)
}

func (s *Store) sendInitialMatchData(u *User, channelName string) error {
	serverSteamId := strings.Split(channelName, ".")[1]
	match, err := api.SteamApi.GetRealTimeStats(serverSteamId)
	if err != nil {
		return err
	}
	for i, game := range match.Teams {
		for playerKey, player := range game.Players {
			match.Teams[i].Players[playerKey].Hero = storage.HeroById(player.HeroId)
		}
	}

	marshal, err := json.Marshal(match)
	if err != nil {
		return err
	}

	var apiMatch dota.ApiLiveMatch
	if err := json.Unmarshal(marshal, &apiMatch); err != nil {
		return err
	}

	log.Infof(`Sending api response to user %s`, u.ID)
	u.send <- &ApiMatchResponse{
		Event:   channelName,
		Data:    apiMatch,
		Success: false,
	}

	return nil
}

func isLiveMatchChannel(channelName string) bool {
	return strings.HasPrefix(channelName, channelLiveMatchPrefix)
}

func (s *Store) findAndDeliver(channel string, content string) {

	var match dota.ApiLiveMatch
	if err := json.Unmarshal([]byte(content), &match); err != nil {
		log.WithFields(log.Fields{
			"fn": "findAndDeliver",
		}).Error(err)
	}

	wsResp := &ApiMatchResponse{
		Event:   channel,
		Data:    match,
		Success: true,
	}

	s.Lock()
	if _, ok := s.Channels[channel]; ok {
		log.Infof("Broadcasting to %s, user count %d \n", channel, len(s.Channels[channel]))
		start := time.Now()
		for _, u := range s.Channels[channel] {
			s.Channels[channel][u.ID].send <- wsResp
		}
		elapsed := time.Since(start)
		log.Infof("Delivered in took %s", elapsed)
	} else {
		log.Errorf("Channel %s not found at our store\n", channel)
	}
	s.Unlock()
}

func (s *Store) removeUser(u *User) {

	s.Lock()
	s.removeUserFromChannels(u)

	for i, storeUser := range s.Users {
		if storeUser == u {
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
			log.Debugf("Removed user %s from store \n", u.ID)
		}
	}
	s.Unlock()
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
