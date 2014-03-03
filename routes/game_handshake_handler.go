package routes

import (
	"code.google.com/p/go.net/websocket"
	"regexp"
)

type Player struct{}

type Game struct {
	gameReadyness chan bool
}

func (g *Game) AddPlayer(p *Player) error {
	return nil
}

func FindOrCreateGameById(id string) *Game {
	return &Game{}
}

func GameHandshakeHandler(ws *websocket.Conn) {
	re := regexp.MustCompile("[^/]+$")
	id := re.FindString(ws.Request().URL.Path)
	game := FindOrCreateGameById(id)

	game.AddPlayer(&Player{})

	for {
		// fired after number of players changes
		isGameReady := <-game.gameReadyness

		if isGameReady {
			websocket.JSON.Send(ws, map[string]string{"handshake": "ready"})
			break
		} else {
			websocket.JSON.Send(ws, map[string]string{"handshake": "wait"})
		}
	}

	ws.Close()
}
