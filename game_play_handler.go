package main

import (
	"errors"
	"net/http"

	sb "github.com/artemave/conways-go/synchronized_broadcaster"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nu7hatch/gouuid"
)

type Player struct {
	GameServerMessages      chan sb.BroadcastMessage
	id                      *uuid.UUID
	SynchronizedBroadcaster *sb.SynchronizedBroadcaster
}

func NewPlayer(game *Game) *Player {
	u4, _ := uuid.NewV4()
	player := &Player{
		id: u4,
		SynchronizedBroadcaster: game.SynchronizedBroadcaster,
		GameServerMessages:      make(chan sb.BroadcastMessage),
	}
	return player
}

func (p Player) ClientId() *uuid.UUID {
	return p.id
}

func (p Player) Inbox() chan sb.BroadcastMessage {
	return p.GameServerMessages
}

func (p Player) MessageAcknowledged() {
	p.SynchronizedBroadcaster.MessageAcknowledged(p)
}

type Game struct {
	Id                      string
	SynchronizedBroadcaster *sb.SynchronizedBroadcaster
	Conway                  *conway.Game
}

func NewGame(id string) *Game {
	game := &Game{
		Id: id,
		SynchronizedBroadcaster: sb.NewSynchronizedBroadcaster(),
		Conway:                  &conway.Game{Cols: 300, Rows: 200},
	}
	return game
}

func (g *Game) AddPlayer() (*Player, error) {
	if len(g.SynchronizedBroadcaster.Clients) >= 2 {
		return &Player{}, errors.New("Game has already reached maximum number players")
	}
	p := NewPlayer(g)

	enoughPlayersToStart := len(g.SynchronizedBroadcaster.Clients) >= 2
	if enoughPlayersToStart {
		g.SynchronizedBroadcaster.SendBroadcastMessage(enoughPlayersToStart)
	}

	return p, nil
}

func (g *Game) NextGeneration() *conway.Generation {
	return g.StartGeneration()
}

func (g *Game) StartGeneration() *conway.Generation {
	// 150x100
	return &conway.Generation{
		{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},

		{Point: conway.Point{Row: 95, Col: 145}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 96, Col: 145}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 96, Col: 146}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 95, Col: 146}, State: conway.Live, Player: conway.Player2},
	}
}

func (g *Game) RemovePlayer(p *Player) error {
	if err := g.SynchronizedBroadcaster.RemoveClient(p); err != nil {
		return err
	}

	g.SynchronizedBroadcaster.SendBroadcastMessage(len(g.SynchronizedBroadcaster.Clients) >= 2)

	return nil
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
		case msg := <-player.GameServerMessages:

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

			player.MessageAcknowledged()
		case <-disconnected:
			return
		}
	}
}
