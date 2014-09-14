package game_result_calculator

import (
	. "github.com/artemave/conways-go/conway"
)

type GameResultCalculator func(*Generation, []*Player, interface {
	WinSpot(*Player) *Point
}) *Player

func (grc GameResultCalculator) Winner(g *Generation, ps []*Player, game interface {
	WinSpot(*Player) *Point
}) *Player {
	return grc(g, ps, game)
}

var CaptureFlagCalculator = GameResultCalculator(
	func(generation *Generation, players []*Player, game interface {
		WinSpot(*Player) *Point
	}) *Player {
		return nil
	},
)
