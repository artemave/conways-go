package main

import (
	"errors"
	"net/http"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	"github.com/artemave/conways-go/dependencies/gouuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type GameServerMessage struct {
	Id   *uuid.UUID
	Data interface{}
}

func NewGameServerMessage(v interface{}) *GameServerMessage {
	u4, _ := uuid.NewV4()
	msg := &GameServerMessage{
		Id:   u4,
		Data: v,
	}
	return msg
}

type Player struct {
	GameServerMessage chan *GameServerMessage
	GameData          chan *conway.Generation
	Id                *uuid.UUID
	Game              *Game
}

func NewPlayer(game *Game) *Player {
	u4, _ := uuid.NewV4()
	player := &Player{
		Id:                u4,
		Game:              game,
		GameServerMessage: make(chan *GameServerMessage),
	}
	return player
}

func (self *Player) MessageAcknowledged(msgId *uuid.UUID) {
	self.Game.MessageAcknowledged <- msgId
}

type Game struct {
	Id                  string
	Players             []*Player
	MessageAcknowledged chan *uuid.UUID
}

func NewGame(id string) *Game {
	game := &Game{
		Id:                  id,
		Players:             []*Player{},
		MessageAcknowledged: make(chan *uuid.UUID),
	}
	go func() {
		<-game.MessageAcknowledged
	}()
	return game
}

func (g *Game) AddPlayer() (*Player, error) {
	if len(g.Players) >= 2 {
		return &Player{}, errors.New("Game has already reached maximum number players")
	}
	p := NewPlayer(g)
	g.Players = append(g.Players, p)

	msg := NewGameServerMessage(len(g.Players) >= 2)

	go func() {
		for _, p := range g.Players {
			p.GameServerMessage <- msg
		}
	}()

	return p, nil
}

func (self *Game) RemovePlayer(p *Player) error {
	for i, player := range self.Players {
		if *player.Id == *p.Id {
			self.Players = append(self.Players[:i], self.Players[i+1:]...)

			msg := NewGameServerMessage(len(self.Players) >= 2)

			go func() {
				for _, p := range self.Players {
					p.GameServerMessage <- msg
				}
			}()

			return nil
		}
	}
	return errors.New("Trying to delete non-existent player")
}

type GamesRepo struct {
	Games []*Game
}

func NewGamesRepo() *GamesRepo {
	gr := &GamesRepo{
		Games: []*Game{},
	}
	return gr
}

// FIXME this is not thread-safe
func (gr *GamesRepo) FindOrCreateGameById(id string) *Game {
	for _, game := range gr.Games {
		if game.Id == id {
			return game
		}
	}
	newGame := NewGame(id)
	gr.Games = append(gr.Games, newGame)
	return newGame
}

var gamesRepo = NewGamesRepo()

func GamePlayHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		gou.Error(err)
		return
	}
	defer ws.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	gou.Debug("/games/%v", id)

	game := gamesRepo.FindOrCreateGameById(id)

	player, err := game.AddPlayer()

	if err != nil {
		ws.WriteJSON(map[string]string{"handshake": "game_taken"})
		return
	} else {
		defer game.RemovePlayer(player)
	}

	disconnected := make(chan bool)

	go func() {
		for {
			if _, _, err := ws.NextReader(); err != nil {
				disconnected <- true
			}
		}
	}()

	for {
		select {
		case msg := <-player.GameServerMessage:

			switch messageData := msg.Data.(type) {
			case bool:
				if messageData {
					if err := ws.WriteJSON(map[string]string{"handshake": "ready"}); err != nil {
						gou.Error("Send to user: ", err)
						return
					}
				} else {
					if err := ws.WriteJSON(map[string]string{"handshake": "wait"}); err != nil {
						gou.Error("Send to user: ", err)
						return
					}
				}
			case *conway.Generation:
				points := conway.GenerationToPoints(messageData)
				if err := ws.WriteJSON(points); err != nil {
					gou.Error("Send to user: ", err)
					return
				}
			}

			player.MessageAcknowledged(msg.Id)
		case <-disconnected:
			return
		}
	}
}
