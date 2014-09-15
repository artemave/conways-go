package game_test

import (
	"time"

	"github.com/artemave/conways-go/conway"
	. "github.com/artemave/conways-go/game"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testDelay = &Delay

type StubBroadcaster struct {
	Message interface{}
}

func (b *StubBroadcaster) Clients() []sb.SynchronizedBroadcasterClient {
	return []sb.SynchronizedBroadcasterClient{}
}
func (b *StubBroadcaster) AddClient(sb.SynchronizedBroadcasterClient)    {}
func (b *StubBroadcaster) RemoveClient(sb.SynchronizedBroadcasterClient) {}
func (b *StubBroadcaster) SendBroadcastMessage(data interface{}) {
	b.Message = data
}
func (b *StubBroadcaster) MessageAcknowledged(sb.SynchronizedBroadcasterClient) {}

type StubResultCalculator func(*conway.Generation, interface {
	WinSpot(*conway.Player) *conway.Point
}) *conway.Player

func (self StubResultCalculator) Winner(generation *conway.Generation, game interface {
	WinSpot(*conway.Player) *conway.Point
}) *conway.Player {
	return self(generation, game)
}

var testGeneration = &conway.Generation{
	{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
}

var _ = Describe("Game", func() {
	*testDelay = time.Duration(30)
	var (
		game    *Game
		player1 *Player
		player2 *Player
	)

	Describe("StartClock", func() {
		BeforeEach(func() {
			game = NewGame("id", "small", testGeneration)
			game.Broadcaster = &StubBroadcaster{}
			player1, _ = game.AddPlayer()
			player2, _ = game.AddPlayer()
		})
		AfterEach(func() {
			game.StopClock()
		})
		Context("Game is on", func() {
			BeforeEach(func() {
				game.GameResultCalculator = StubResultCalculator(
					func(generation *conway.Generation, game interface {
						WinSpot(*conway.Player) *conway.Point
					}) *conway.Player {
						return nil
					},
				)
				game.StartClock()
			})
			It("broadcasts next generation", func() {
				Eventually(func() interface{} {
					stubBroadcaster := game.Broadcaster.(*StubBroadcaster)
					return stubBroadcaster.Message
				}).Should(BeAssignableToTypeOf(&conway.Generation{}))
			})
		})
		Context("Game is complete", func() {
			BeforeEach(func() {
				game.GameResultCalculator = StubResultCalculator(
					func(generation *conway.Generation, game interface {
						WinSpot(*conway.Player) *conway.Point
					}) *conway.Player {
						r := conway.Player1
						return &r
					},
				)
				game.StartClock()
			})
			It("broadcasts game result", func() {
				Eventually(func() interface{} {
					stubBroadcaster := game.Broadcaster.(*StubBroadcaster)
					return stubBroadcaster.Message
				}).Should(BeAssignableToTypeOf(GameResult{}))
			})
		})
	})
})
