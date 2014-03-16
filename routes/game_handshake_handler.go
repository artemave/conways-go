package routes

import (
  "errors"
  "github.com/araddon/gou"
  "github.com/gorilla/websocket"
  "net/http"
  "regexp"
)

type Player struct {
  GameReady chan bool
}

func NewPlayer() *Player {
  player := &Player{
    GameReady: make(chan bool),
  }
  return player
}

type Game struct {
  Id      string
  Players []*Player
}

func NewGame(id string) *Game {
  game := &Game{
    Id:      id,
    Players: []*Player{},
  }
  return game
}

func (g *Game) AddPlayer(p *Player) error {
  if len(g.Players) >= 2 {
    return errors.New("Game has already reached maximum number players")
  }

  g.Players = append(g.Players, p)
  go func() {
    for _, p := range g.Players {
      if len(g.Players) >= 2 {
        p.GameReady <- true
      } else {
        p.GameReady <- false
      }
    }
  }()
  return nil
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
      found := false

      for _, game := range gr.Games {
        if game.Id == id {
          gr.GetGameResult <- game
          found = true
          break
        }
      }

      if !found {
        newGame := NewGame(id)
        gr.Games = append(gr.Games, newGame)
        gr.GetGameResult <- newGame
      }
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

func GameHandshakeHandler(w http.ResponseWriter, r *http.Request) {
  ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
  if _, ok := err.(websocket.HandshakeError); ok {
    http.Error(w, "Not a websocket handshake", 400)
    return
  } else if err != nil {
    gou.Error(err)
    return
  }
  defer ws.Close()

  re := regexp.MustCompile("[^/]+$")
  id := re.FindString(r.URL.Path)
  game := gamesRepo.FindOrCreateGameById(id)

  player := NewPlayer()
  err = game.AddPlayer(player)

  if err != nil {
    ws.WriteJSON(map[string]string{"handshake": "game_taken"})
    return
  }

  for {
    // fired after number of players changes
    if <-player.GameReady {
      ws.WriteJSON(map[string]string{"handshake": "ready"})
    } else {
      ws.WriteJSON(map[string]string{"handshake": "wait"})
    }
  }
}
