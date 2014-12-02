package game

import (
	"errors"
	"time"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	grc "github.com/artemave/conways-go/game_result_calculator"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
)

var Delay = time.Duration(1000)

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
		Id:                   id,
		Broadcaster:          sb.NewSynchronizedBroadcaster(),
		GameResultCalculator: grc.CaptureFlagCalculator,
		Conway:               &conway.Game{Cols: cols, Rows: rows},
		stopClock:            make(chan bool, 1),
		clientCells:          make(chan []conway.Cell),
		players:              []*Player{},
		startGeneration:      startGeneration,
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
	for _, p := range g.players {
		if p.ClientId() != player.ClientId() {
			p.GamePauseMessages <- true
		}
	}
}

func (g *Game) ResumeBy(player *Player) {
	for _, p := range g.players {
		if p.ClientId() != player.ClientId() {
			p.GamePauseMessages <- false
		}
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
				if winnerIndex := g.GameResultCalculator.Winner(g.currentGeneration, g); winnerIndex != nil {
					g.Broadcaster.SendBroadcastMessage(GameResult{g.playerByIndex(winnerIndex)})
					return
				} else {
					g.Broadcaster.SendBroadcastMessage(g.NextGeneration())
				}
				time.Sleep(Delay * time.Millisecond)
			}
		}
	}()
}

func (g *Game) playerIndexes() []*conway.Player {
	playerIndexes := []*conway.Player{}
	for _, v := range g.players {
		playerIndexes = append(playerIndexes, &v.PlayerIndex)
	}
	return playerIndexes
}

func (self *Game) playerByIndex(idx *conway.Player) *Player {
	for _, v := range self.players {
		if v.PlayerIndex == *idx {
			return v
		}
	}
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
	gou.Debug("Removing player ", p.id)

	g.Broadcaster.RemoveClient(p)
	p.CleanUp()

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

	gou.Debug("Player removed ", p.id)
	return nil
}
