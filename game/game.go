package game

import (
	"errors"
	"time"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/clock"
	"github.com/artemave/conways-go/conway"
	grc "github.com/artemave/conways-go/game_result_calculator"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
)

var Delay = time.Duration(1000)

type PauseGame bool
type PlayersAreReady bool

type GameResult struct {
	Winner *Player
}

type WinSpot struct {
	Player conway.Player
	Point  conway.Point
}

type Broadcaster interface {
	Clients() []sb.SynchronizedBroadcasterClient
	AddClient(sb.SynchronizedBroadcasterClient)
	RemoveClient(sb.SynchronizedBroadcasterClient)
	SendBroadcastMessage(interface{})
	MessageAcknowledged(sb.SynchronizedBroadcasterClient)
}

type Game struct {
	Id     string
	Conway *conway.Game
	Broadcaster
	GameResultCalculator interface {
		Winner(*conway.Generation, interface {
			WinSpot(*conway.Player) *conway.Point
		}) *conway.Player
	}
	currentGeneration *conway.Generation
	startGeneration   *conway.Generation
	Players           []*Player
	clientCells       chan []conway.Cell
	PausedByPlayer    conway.Player
	IsPractice        bool
	clock             *clock.Clock
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
		Id:                   id,
		Broadcaster:          sb.NewSynchronizedBroadcaster(),
		GameResultCalculator: grc.CaptureFlagCalculator,
		Conway:               &conway.Game{Cols: cols, Rows: rows},
		clientCells:          make(chan []conway.Cell),
		Players:              []*Player{},
		startGeneration:      startGeneration,
		PausedByPlayer:       conway.None,
		IsPractice:           false,
		clock:                clock.NewClock(Delay),
	}

	// TODO clean up clock when game is destroyed
	// TODO clean this up when game is destroyed
	go func() {
		for {
			select {
			case cells := <-game.clientCells:
				game.currentGeneration.AddCells(cells)
			case <-game.clock.NextTick():
				winnerIndex := game.GameResultCalculator.Winner(game.currentGeneration, game)

				if winnerIndex != nil {
					game.Broadcaster.SendBroadcastMessage(GameResult{game.playerByIndex(winnerIndex)})
					game.StopClock()
				} else {
					game.Broadcaster.SendBroadcastMessage(game.NextGeneration())
				}
			}
		}
	}()

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
	g.Players = append(g.Players, p)
	gou.Debug("Started adding a player ", p.id)

	g.Broadcaster.AddClient(p)

	enoughPlayersToStart := PlayersAreReady(len(g.Broadcaster.Clients()) >= 2)
	g.Broadcaster.SendBroadcastMessage(enoughPlayersToStart)

	pNum := conway.Player1
	if enoughPlayersToStart {
		// TODO test player number assignment (when client reconnects)

		for _, p := range g.Players {
			if p.PlayerIndex == conway.Player1 {
				pNum = conway.Player2
			}
		}

		g.StartClock()
	}
	p.PlayerIndex = pNum

	gou.Debug("Player added ", p.id)
	return p, nil
}

func (g *Game) WinSpot(playerIndex *conway.Player) *conway.Point {
	if *playerIndex == conway.Player1 {
		return &conway.Point{Col: g.Conway.Cols - 3, Row: g.Conway.Rows - 3}
	} else {
		return &conway.Point{Col: 2, Row: 2}
	}
}

func (g *Game) WinSpots() []WinSpot {
	winSpots := []WinSpot{}
	for _, playerIndex := range g.playerIndexes() {
		winSpots = append(winSpots, WinSpot{Player: *playerIndex, Point: *g.WinSpot(playerIndex)})
	}
	return winSpots
}

func (g *Game) PauseBy(player *Player) {
	g.StopClock()
	g.PausedByPlayer = player.PlayerIndex
	g.Broadcaster.SendBroadcastMessage(PauseGame(true))
}

func (g *Game) Resume() {
	g.Broadcaster.SendBroadcastMessage(PauseGame(false))
	g.PausedByPlayer = conway.None
	g.StartClock()
}

func (g *Game) IsPaused() bool {
	return g.PausedByPlayer != conway.None
}

func (g *Game) StartClock() {
	g.clock.StartClock()
}

func (g *Game) playerIndexes() []*conway.Player {
	playerIndexes := []*conway.Player{}
	for _, v := range g.Players {
		playerIndexes = append(playerIndexes, &v.PlayerIndex)
	}
	return playerIndexes
}

func (self *Game) playerByIndex(idx *conway.Player) *Player {
	for _, v := range self.Players {
		if v.PlayerIndex == *idx {
			return v
		}
	}
	return nil
}

func (g *Game) StopClock() {
	g.clock.StopClock()
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
	gou.Debug("Removing player ", p.id)

	g.Broadcaster.RemoveClient(p)
	p.CleanUp()

	newPlayers := []*Player{}
	for _, pl := range g.Players {
		if pl.id != p.id {
			newPlayers = append(newPlayers, pl)
		}
	}
	g.Players = newPlayers

	enoughPlayersToStart := PlayersAreReady(len(g.Broadcaster.Clients()) >= 2)
	g.Broadcaster.SendBroadcastMessage(enoughPlayersToStart)

	if !enoughPlayersToStart {
		g.StopClock()
	}

	gou.Debug("Player removed ", p.id)
	return nil
}
