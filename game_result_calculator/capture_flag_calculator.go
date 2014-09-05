package game_result_calculator

import (
	. "github.com/artemave/conways-go/conway"
)

type ResultCalculator func(*Generation, []*Player) *Player

func (rc ResultCalculator) Winner(generation *Generation, players []*Player) *Player {
	return rc(generation, players)
}

var CaptureFlagCalculator = ResultCalculator(
	func(generation *Generation, players []*Player) *Player {
		return nil
	},
)
