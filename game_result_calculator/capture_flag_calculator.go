package game_result_calculator

import (
	. "github.com/artemave/conways-go/conway"
)

type GameResultCalculator func(*Generation, interface {
	WinSpot(*Player) *Point
}) *Player

func (grc GameResultCalculator) Winner(g *Generation, game interface {
	WinSpot(*Player) *Point
}) *Player {
	return grc(g, game)
}

var CaptureFlagCalculator = GameResultCalculator(
	func(generation *Generation, game interface {
		WinSpot(*Player) *Point
	}) *Player {
		if generation == nil {
			return nil
		}

		draw := None
		var winner *Player
		playerHasLiveCells := make(map[Player]bool)

		for _, cell := range *generation {
			playerHasLiveCells[cell.Player] = true

			if cell.Point == *game.WinSpot(&cell.Player) {
				if winner != nil {
					return &draw
				}
				winner = &cell.Player
			}
		}

		if len(playerHasLiveCells) == 1 {
			for k, _ := range playerHasLiveCells {
				return &k
			}
		}

		if len(playerHasLiveCells) == 0 {
			return &draw
		}

		return winner
	},
)
