package routes

import (
	"code.google.com/p/go.net/websocket"
	"regexp"
)

type Player struct{}

type Game struct {
	Id        string
	GameReady chan bool
	Players   []*Player
}

func NewGame(id string) *Game {
	game := &Game{
		Id:        id,
		GameReady: make(chan bool),
		Players:   []*Player{},
	}
	return game
}

func (g *Game) AddPlayer(p *Player) {
	g.Players = append(g.Players, p)
	go func() {
		if len(g.Players) >= 2 {
			g.GameReady <- true
		} else {
			g.GameReady <- false
		}
	}()
}

type GamesRepo struct {
	GetGame       chan string
	GetGameResult chan *Game
	Games         []*Game
}

func NewGamesRepo() *GamesRepo {
	gr := &GamesRepo{
		Games:         []*Game{},
		GetGame:       make(chan string),
		GetGameResult: make(chan *Game),
	}

	go func() {
		for {
			id := <-gr.GetGame
			for i := 0; i < len(gr.Games); i++ {
				if gr.Games[i].Id == id {
					gr.GetGameResult <- gr.Games[i]
					break
				}
			}
			newGame := NewGame(id)
			gr.Games = append(gr.Games, newGame)
			gr.GetGameResult <- newGame
		}
	}()
	return gr
}

func (gr *GamesRepo) FindOrCreateGameById(id string) *Game {
	gr.GetGame <- id
	game := <-gr.GetGameResult
	return game
}

var gamesRepo = NewGamesRepo()

func GameHandshakeHandler(ws *websocket.Conn) {
	re := regexp.MustCompile("[^/]+$")
	id := re.FindString(ws.Request().URL.Path)
	game := gamesRepo.FindOrCreateGameById(id)

	game.AddPlayer(&Player{})

	for {
		// fired after number of players changes
		if <-game.GameReady {
			websocket.JSON.Send(ws, map[string]string{"handshake": "ready"})
			break
		} else {
			websocket.JSON.Send(ws, map[string]string{"handshake": "wait"})
		}
	}

	ws.Close()
}
