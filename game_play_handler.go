package main

import (
	"errors"
	"net/http"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/dependencies/gouuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Player struct {
	GameReady chan bool
	Id        *uuid.UUID
}

func NewPlayer() *Player {
	u4, _ := uuid.NewV4()
	player := &Player{
		GameReady: make(chan bool),
		Id:        u4,
	}
	return player
}

type Game struct {
	Id      string
	Players []*Player
}

func NewGame(id string) *Game {
	game := &Game{
		Id:      id,
		Players: []*Player{},
	}
	return game
}

func (g *Game) AddPlayer(p *Player) error {
	if len(g.Players) >= 2 {
		return errors.New("Game has already reached maximum number players")
	}
	g.Players = append(g.Players, p)
	go func() {
		for _, p := range g.Players {
			if len(g.Players) >= 2 {
				p.GameReady <- true
			} else {
				p.GameReady <- false
			}
		}
	}()
	return nil
}

func (self *Game) RemovePlayer(p *Player) error {
	for i, player := range self.Players {
		if *player.Id == *p.Id {
			self.Players = append(self.Players[:i], self.Players[i+1:]...)

			go func() {
				for _, p := range self.Players {
					if len(self.Players) >= 2 {
						p.GameReady <- true
					} else {
						p.GameReady <- false
					}
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

	player := NewPlayer()
	err = game.AddPlayer(player)

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
		case ready := <-player.GameReady:
			// fired after number of players changes
			if ready {
				ws.WriteJSON(map[string]string{"handshake": "ready"})
			} else {
				ws.WriteJSON(map[string]string{"handshake": "wait"})
			}
		case <-disconnected:
			return
		}
	}
}
