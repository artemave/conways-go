package game

import (
	"fmt"

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
				msg, ok := <-player.GameServerMessages

				if !ok {
					fmt.Printf("NOT OK\n")
					return
				}

				fmt.Printf("DUMMY PLAYER ACK %#v\n", msg)

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
