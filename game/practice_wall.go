package game

import (
	"github.com/araddon/gou"
)

type PracticeWall struct {
	*Game
	dummyPlayer *Player
}

func NewPracticeWall(game *Game) *PracticeWall {
	pw := &PracticeWall{
		Game: game,
	}

	go func() {
		player, err := game.AddPlayer()
		if err != nil {
			gou.Error("Failed to add practice player: ", err)
		}
		pw.dummyPlayer = player

		go func() {
			for {
				_, ok := <-player.GameServerMessages

				if !ok {
					return
				}

				player.MessageAcknowledged()
			}
		}()
	}()

	return pw
}

func (pw *PracticeWall) RemoveDummyPlayer() {
	if pw.dummyPlayer != nil {
		pw.Game.RemovePlayer(pw.dummyPlayer)
	}
}
