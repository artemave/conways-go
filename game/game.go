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

const maxFreeCells = 10
const restoreFreeCellsEveryNTicks = 2

type PauseGame bool
type PlayersAreReady bool
type CellCount int

type NewCellsCache struct {
	Cells          []conway.Cell
	FreeCellsCount CellCount
}

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
	size   string
	Conway *conway.Game
	Broadcaster
	GameResultCalculator interface {
		Winner(*conway.Generation, interface {
			WinSpot(*conway.Player) *conway.Point
		}) *conway.Player
	}
	currentGeneration *conway.Generation
	Players           []*Player
	clientCells       chan []conway.Cell
	PausedByPlayer    conway.Player
	isPractice        bool
	isFinished        bool
	scoredBy          string
	clock             *clock.Clock
	newCellsCache     map[conway.Player]*NewCellsCache
}

type GameTicker struct {
	ticker time.Ticker
}

func (gt GameTicker) C() <-chan time.Time {
	return gt.ticker.C
}

func (gt GameTicker) Stop() {
	gt.ticker.Stop()
}

func tickerFactory(delay time.Duration) clock.Ticker {
	return GameTicker{ticker: *time.NewTicker(delay)}
}

func NewGame(id string, size string, startGeneration *conway.Generation, isPractice bool) *Game {
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
		size:                 size,
		Broadcaster:          sb.NewSynchronizedBroadcaster(),
		GameResultCalculator: grc.CaptureFlagCalculator,
		Conway:               &conway.Game{Cols: cols, Rows: rows},
		clientCells:          make(chan []conway.Cell),
		Players:              []*Player{},
		currentGeneration:    startGeneration,
		PausedByPlayer:       conway.None,
		isPractice:           isPractice,
		clock:                clock.NewClock(Delay, tickerFactory),
		newCellsCache: map[conway.Player]*NewCellsCache{
			conway.Player1: &NewCellsCache{
				FreeCellsCount: CellCount(restoreFreeCellsEveryNTicks * maxFreeCells),
				Cells:          []conway.Cell{},
			},
			conway.Player2: &NewCellsCache{
				FreeCellsCount: CellCount(restoreFreeCellsEveryNTicks * maxFreeCells),
				Cells:          []conway.Cell{},
			},
		},
	}

	// TODO clean up clock when game is destroyed
	// TODO clean this up when game is destroyed
	go func() {
		for {
			select {
			case cells := <-game.clientCells:
				if len(cells) > 0 {
					c := game.newCellsCache[cells[0].Player]

					c.FreeCellsCount -= CellCount(restoreFreeCellsEveryNTicks * len(cells))
					c.Cells = append(c.Cells, cells...)
				}
			case <-game.clock.NextTick():
				winnerIndex := game.GameResultCalculator.Winner(game.currentGeneration, game)

				if winnerIndex != nil {
					game.Broadcaster.SendBroadcastMessage(GameResult{game.playerByIndex(winnerIndex)})
					game.isFinished = true
					game.StopClock()
				} else {
					nextGeneration := game.NextGeneration()

					for pi, _ := range game.newCellsCache {
						c := game.newCellsCache[pi]
						if c.FreeCellsCount/restoreFreeCellsEveryNTicks < maxFreeCells {
							c.FreeCellsCount += 1
						}

						if len(c.Cells) > 0 {
							nextGeneration.AddCells(c.Cells)
							c.Cells = []conway.Cell{}
						}
					}

					game.Broadcaster.SendBroadcastMessage(nextGeneration)
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

	// TODO test player number assignment (when client reconnects)
	pNum := conway.Player1
	for _, p := range g.Players {
		if p.PlayerIndex == conway.Player1 {
			pNum = conway.Player2
		}
	}

	p := NewPlayer(g.Broadcaster, pNum)
	g.Players = append(g.Players, p)
	gou.Debug("Started adding a player ", p.id)

	g.Broadcaster.AddClient(p)

	enoughPlayersToStart := PlayersAreReady(len(g.Broadcaster.Clients()) >= 2)
	g.Broadcaster.SendBroadcastMessage(enoughPlayersToStart)

	if enoughPlayersToStart {
		g.StartClock()
	}

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
	// so that players get served initial game state immediately
	// after game starts (without having to wait for next clock tick)
	go g.Broadcaster.SendBroadcastMessage(g.currentGeneration)
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
	g.currentGeneration = g.Conway.NextGeneration(g.currentGeneration)
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

func (g *Game) FreeCellsCountOf(player *Player) CellCount {
	return g.newCellsCache[player.PlayerIndex].FreeCellsCount / restoreFreeCellsEveryNTicks
}

// SetScoredBy : mark game scored, so that no more scores can be submitted for this game
func (g *Game) SetScoredBy(googlePlayerID string) {
	g.scoredBy = googlePlayerID
}

// GetScoredBy : check if the game score was already submitted
func (g *Game) GetScoredBy() string {
	return g.scoredBy
}

// Size : game size ("small", "medium", "large")
func (g *Game) Size() string {
	return g.size
}

// IsPractice : is this a practice wall or a duel between humans?
func (g *Game) IsPractice() bool {
	return g.isPractice
}

// IsFinished : has this game already finished?
func (g *Game) IsFinished() bool {
	return g.isFinished
}
