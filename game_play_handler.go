package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	sb "github.com/artemave/conways-go/synchronized_broadcaster"

	"code.google.com/p/go-uuid/uuid"
	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var delay = time.Duration(1000)

type Player struct {
	GameServerMessages      chan sb.BroadcastMessage
	id                      string
	SynchronizedBroadcaster *sb.SynchronizedBroadcaster
}

func NewPlayer(game *Game) *Player {
	u4 := uuid.New()
	player := &Player{
		id: u4,
		SynchronizedBroadcaster: game.SynchronizedBroadcaster,
		GameServerMessages:      make(chan sb.BroadcastMessage),
	}
	return player
}

func (p Player) ClientId() string {
	return p.id
}

func (p Player) Inbox() chan sb.BroadcastMessage {
	return p.GameServerMessages
}

func (p Player) MessageAcknowledged() {
	p.SynchronizedBroadcaster.MessageAcknowledged()
}

type Game struct {
	Id                      string
	SynchronizedBroadcaster *sb.SynchronizedBroadcaster
	Conway                  *conway.Game
	currentGeneration       *conway.Generation
	stopClock               chan bool
	playerNumbers           map[string]int
}

func NewGame(id string) *Game {
	game := &Game{
		Id: id,
		SynchronizedBroadcaster: sb.NewSynchronizedBroadcaster(),
		Conway:                  &conway.Game{Cols: 300, Rows: 200},
		stopClock:               make(chan bool, 1),
		playerNumbers:           make(map[string]int),
	}
	return game
}

func (g *Game) AddPlayer() (*Player, error) {
	if len(g.SynchronizedBroadcaster.Clients) >= 2 {
		return &Player{}, errors.New("Game has already reached maximum number players")
	}
	p := NewPlayer(g)

	g.SynchronizedBroadcaster.AddClient(p)

	enoughPlayersToStart := len(g.SynchronizedBroadcaster.Clients) >= 2
	g.SynchronizedBroadcaster.SendBroadcastMessage(enoughPlayersToStart)

	pNum := 1
	if enoughPlayersToStart {
		// TODO test player number assignment (when client reconnects)

		// there is only one element in playerNumbers at this point
		for _, v := range g.playerNumbers {
			if v == 1 {
				pNum = 2
			}
		}
		g.StartClock()
	}
	g.playerNumbers[p.id] = pNum

	return p, nil
}

func (g *Game) PlayerNumber(player *Player) int {
	return g.playerNumbers[player.id]
}

func (g *Game) StartClock() {
	go func() {
		for {
			select {
			case <-g.stopClock:
				return
			default:
				g.SynchronizedBroadcaster.SendBroadcastMessage(g.NextGeneration())
				time.Sleep(delay * time.Millisecond)
			}
		}
	}()
}

func (g *Game) StopClock() {
	g.stopClock <- true
}

func (g *Game) NextGeneration() *conway.Generation {
	if g.currentGeneration == nil {
		g.currentGeneration = startGeneration
	} else {
		g.currentGeneration = g.Conway.NextGeneration(g.currentGeneration)
	}
	return g.currentGeneration
}

// 150x100
var startGeneration = &conway.Generation{
	{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 4, Col: 3}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},

	{Point: conway.Point{Row: 95, Col: 145}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 96, Col: 145}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 96, Col: 146}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 95, Col: 146}, State: conway.Live, Player: conway.Player2},
}

func (g *Game) RemovePlayer(p *Player) error {
	close(p.GameServerMessages)

	if err := g.SynchronizedBroadcaster.RemoveClient(p); err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	enoughPlayersToStart := len(g.SynchronizedBroadcaster.Clients) >= 2
	g.SynchronizedBroadcaster.SendBroadcastMessage(enoughPlayersToStart)

	if !enoughPlayersToStart {
		g.StopClock()
	}

	delete(g.playerNumbers, p.id)

	gou.Debug("Removing player ", p.id)
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
			var msg map[string]interface{}
			if err := ws.ReadJSON(&msg); err != nil {
				disconnected <- true
				return
			} else {
				switch msg["acknowledged"].(string) {
				case "ready", "wait", "game":
					player.MessageAcknowledged()
				default:
					panic("Unknown client message")
				}
			}
		}
	}()

	for {
		select {
		case msg := <-player.GameServerMessages:

			switch messageData := msg.Data.(type) {
			case bool:
				if messageData {
					if err := ws.WriteJSON(map[string]string{"handshake": "ready", "player": strconv.Itoa(game.PlayerNumber(player))}); err != nil {
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
				if err := ws.WriteJSON(messageData); err != nil {
					gou.Error("Send to user: ", err)
					return
				}
			}
		case <-disconnected:
			fmt.Printf("Client disconnected\n")
			ws.Close()
			return
		}
	}
}
