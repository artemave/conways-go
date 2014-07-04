package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
)

type Game struct {
	Id                      string
	SynchronizedBroadcaster *sb.SynchronizedBroadcaster
	Conway                  *conway.Game
	currentGeneration       *conway.Generation
	stopClock               chan bool
	playerNumbers           map[string]conway.Player
	clientCells             chan []conway.Cell
}

func NewGame(id string) *Game {
	game := &Game{
		Id: id,
		SynchronizedBroadcaster: sb.NewSynchronizedBroadcaster(),
		Conway:                  &conway.Game{Cols: 150, Rows: 100},
		stopClock:               make(chan bool, 1),
		clientCells:             make(chan []conway.Cell),
		playerNumbers:           make(map[string]conway.Player),
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

	pNum := conway.Player1
	if enoughPlayersToStart {
		// TODO test player number assignment (when client reconnects)

		// there is only one element in playerNumbers at this point
		for _, v := range g.playerNumbers {
			if v == conway.Player1 {
				pNum = conway.Player2
			}
		}
		g.StartClock()
	}
	g.playerNumbers[p.id] = pNum

	return p, nil
}

func (g *Game) PlayerNumber(player *Player) int {
	return int(g.playerNumbers[player.id])
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

func (g *Game) AddCells(cells []conway.Cell) {
	go func() {
		g.clientCells <- cells
	}()
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
