package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
)

type Broadcaster interface {
	Clients() []sb.SynchronizedBroadcasterClient
	AddClient(sb.SynchronizedBroadcasterClient)
	RemoveClient(sb.SynchronizedBroadcasterClient) error
	SendBroadcastMessage(interface{})
	MessageAcknowledged()
}

type Game struct {
	Id     string
	Conway *conway.Game
	Broadcaster
	currentGeneration *conway.Generation
	startGeneration   *conway.Generation
	stopClock         chan bool
	players           []*Player
	clientCells       chan []conway.Cell
}

func NewGame(id string, size string, startGeneration *conway.Generation) *Game {
	var cols int
	var rows int

	switch size {
	case "small":
		cols = 40
		rows = 26
	case "medium":
		cols = 80
		rows = 50
	case "large":
		cols = 100
		rows = 70
	}

	game := &Game{
		Id:              id,
		Broadcaster:     sb.NewSynchronizedBroadcaster(),
		Conway:          &conway.Game{Cols: cols, Rows: rows},
		stopClock:       make(chan bool, 1),
		clientCells:     make(chan []conway.Cell),
		players:         []*Player{},
		startGeneration: startGeneration,
	}
	return game
}

func (g *Game) Cols() int {
	return g.Conway.Cols
}

func (g *Game) Rows() int {
	return g.Conway.Rows
}

func (g *Game) AddPlayer() (*Player, error) {
	if len(g.Broadcaster.Clients()) >= 2 {
		return &Player{}, errors.New("Game has already reached maximum number players")
	}
	p := NewPlayer(g)
	g.players = append(g.players, p)

	g.Broadcaster.AddClient(p)

	enoughPlayersToStart := len(g.Broadcaster.Clients()) >= 2
	g.Broadcaster.SendBroadcastMessage(enoughPlayersToStart)

	pNum := conway.Player1
	if enoughPlayersToStart {
		// TODO test player number assignment (when client reconnects)

		for _, p := range g.players {
			if p.PlayerIndex == conway.Player1 {
				pNum = conway.Player2
			}
		}
		g.StartClock()
	}
	p.PlayerIndex = pNum

	return p, nil
}

func (g *Game) WinSpot(playerIndex conway.Player) conway.Point {
	if playerIndex == conway.Player1 {
		return conway.Point{Col: g.Conway.Cols - 3, Row: g.Conway.Rows - 3}
	} else {
		return conway.Point{Col: 3, Row: 3}
	}
}

func (g *Game) StartClock() {
	go func() {
		for {
			select {
			case <-g.stopClock:
				return
			case cells := <-g.clientCells:
				g.currentGeneration.AddCells(cells)
			default:
				if gameResult := g.CalculateGameResult(); gameResult != nil {
					g.Broadcaster.SendBroadcastMessage(gameResult)
					return
				} else {
					g.Broadcaster.SendBroadcastMessage(g.NextGeneration())
				}
				time.Sleep(delay * time.Millisecond)
			}
		}
	}()
}

type GameResult struct {
	Winner Player
}

func (g *Game) CalculateGameResult() *GameResult {
	if g.currentGeneration == nil {
		return nil
	}
	for _, cell := range *g.currentGeneration {
		// player1 wins
		if cell.Player == conway.Player1 && cell.Point == g.WinSpot(conway.Player1) {

		}
		// player2 wins
		if cell.Player == conway.Player2 && cell.Point == g.WinSpot(conway.Player2) {

		}
		// handle draw maybe
	}
	// TODO no live cells for player1 or 2
	return nil
}

func (g *Game) StopClock() {
	g.stopClock <- true
}

func (g *Game) NextGeneration() *conway.Generation {
	if g.currentGeneration == nil {
		g.currentGeneration = g.startGeneration
	} else {
		g.currentGeneration = g.Conway.NextGeneration(g.currentGeneration)
	}
	return g.currentGeneration
}

func (g *Game) AddCells(cells []conway.Cell) {
	go func() {
		g.clientCells <- cells
	}()
}

func (g *Game) RemovePlayer(p *Player) error {
	close(p.GameServerMessages)

	if err := g.Broadcaster.RemoveClient(p); err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	enoughPlayersToStart := len(g.Broadcaster.Clients()) >= 2
	g.Broadcaster.SendBroadcastMessage(enoughPlayersToStart)

	if !enoughPlayersToStart {
		g.StopClock()
	}

	newPlayers := []*Player{}
	for _, pl := range g.players {
		if pl.id != p.id {
			newPlayers = append(newPlayers, pl)
		}
	}
	g.players = newPlayers

	gou.Debug("Removing player ", p.id)
	return nil
}
