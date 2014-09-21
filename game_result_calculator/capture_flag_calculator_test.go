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
	var nonePlayer = conway.None
	var generation = &conway.Generation{}

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
		BeforeEach(func() {
			generation = &conway.Generation{
				{Point: conway.Point{Row: 0, Col: 1}, State: conway.Live, Player: conway.Player1},
				{Point: conway.Point{Row: 1, Col: 0}, State: conway.Live, Player: conway.Player2},
			}
		})
		It("Returns nil", func() {
			p := CaptureFlagCalculator(generation, game)
			Expect(p).To(BeNil())
		})
	})
	Context("Player reaches win spot", func() {
		BeforeEach(func() {
			generation = &conway.Generation{
				{Point: conway.Point{Row: 1, Col: 1}, State: conway.Live, Player: conway.Player1},
				{Point: conway.Point{Row: 0, Col: 1}, State: conway.Live, Player: conway.Player2},
				{Point: conway.Point{Row: 1, Col: 0}, State: conway.Live, Player: conway.Player2},
				{Point: conway.Point{Row: 0, Col: 0}, State: conway.Live, Player: conway.Player1},
			}
		})
		It("Declares that player a winner", func() {
			p := CaptureFlagCalculator(generation, game)
			Expect(p).NotTo(BeNil())
			Expect(*p).To(Equal(conway.Player1))
		})

		Context("Another player reaches win spot at the same time", func() {
			BeforeEach(func() {
				generation = &conway.Generation{
					{Point: conway.Point{Row: 1, Col: 1}, State: conway.Live, Player: conway.Player1},
					{Point: conway.Point{Row: 0, Col: 1}, State: conway.Live, Player: conway.Player2},
					{Point: conway.Point{Row: 1, Col: 0}, State: conway.Live, Player: conway.Player2},
					{Point: conway.Point{Row: 0, Col: 0}, State: conway.Live, Player: conway.Player2},
				}
			})
			It("Declares a draw", func() {
				p := CaptureFlagCalculator(generation, game)
				Expect(p).NotTo(BeNil())
				Expect(*p).To(Equal(nonePlayer))
			})
		})
	})
	Context("Player has no more live cells", func() {
		BeforeEach(func() {
			generation = &conway.Generation{
				{Point: conway.Point{Row: 0, Col: 1}, State: conway.Live, Player: conway.Player2},
				{Point: conway.Point{Row: 1, Col: 0}, State: conway.Live, Player: conway.Player2},
				{Point: conway.Point{Row: 1, Col: 1}, State: conway.Live, Player: conway.None},
			}
		})
		It("Declares the other player a winner", func() {
			p := CaptureFlagCalculator(generation, game)
			Expect(p).NotTo(BeNil())
			Expect(*p).To(Equal(conway.Player2))
		})

		Context("Another player runs out of live cells at the same time", func() {
			BeforeEach(func() {
				generation = &conway.Generation{}
			})
			It("Declares a draw", func() {
				p := CaptureFlagCalculator(generation, game)
				Expect(p).NotTo(BeNil())
				Expect(*p).To(Equal(nonePlayer))
			})
		})
	})
})
