package game_result_calculator_test

import (
	"github.com/artemave/conways-go/conway"
	. "github.com/artemave/conways-go/game_result_calculator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestGame func(*conway.Player) *conway.Point

func (self TestGame) WinSpot(p *conway.Player) *conway.Point {
	return self(p)
}

var _ = Describe("CaptureFlagCalculator", func() {
	var player1, player2 *conway.Player
	var generation *conway.Generation
	var game = TestGame(
		func(p *conway.Player) *conway.Point {
			if *p == conway.Player1 {
				return &conway.Point{1, 1}
			} else {
				return &conway.Point{0, 0}
			}
		},
	)

	Context("No players reached win spot", func() {
		It("Returns nil", func() {
			p := CaptureFlagCalculator(generation, []*conway.Player{player1, player2}, game)
			Expect(p).To(BeNil())
		})
	})
	Context("Player reaches win spot", func() {
		PIt("Declares that player a winner", func() {
			// p := CaptureFlagCalculator(generation, players, game)
			// Expect(p).To(Equal(player1))
		})

		Context("Another player reaches win spot at the same time", func() {
			PIt("Declares a draw")
		})
	})
	Context("Player has no more live cells", func() {
		PIt("Declares the other player a winner")

		Context("Another player runs out of live cells at the same time", func() {
			PIt("Declares a draw")
		})
	})
})
