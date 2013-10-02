package comm

import (
  "github.com/artemave/conways-go/game"
  "time"
)

type Server struct {
  addClient chan *Client
  delClient chan *Client
  clients map[int]*Client
  game game.Game
  extraPoints []game.Point
  pointsFromUsers chan []game.Point
  flushExtraPoints chan bool
}

func NewServer() *Server {
  s := &Server{
    addClient: make(chan *Client),
    delClient: make(chan *Client),
    clients: make(map[int]*Client),
    extraPoints: []game.Point{},
    flushExtraPoints: make(chan bool),
    pointsFromUsers: make(chan []game.Point),
    game: game.Game{Rows: 3000, Cols: 4000},
  }
  go s.ListenClientEvents()
  return s
}

func (this *Server) ListenClientEvents() {
  for {
    select {
    case client := <-this.addClient:
      this.clients[client.id] = client
    case client := <-this.delClient:
      delete(this.clients, client.id)
    case points := <-this.pointsFromUsers:
      for _, point := range points {
        this.extraPoints = append(this.extraPoints, point)
      }
    case <-this.flushExtraPoints:
      this.extraPoints = []game.Point{}
    }
  }
}

func (this *Server) ServeTheGame() {
  current_generation := game.GosperGliderGun()
  for {
    time.Sleep(time.Millisecond * 300)

    if len(this.clients) > 0 {
      next_generation := this.game.NextGeneration(current_generation)

      next_generation.AddPoints(this.extraPoints)
      this.flushExtraPoints <- true

      for _, client := range this.clients {
        client.outToUser <- next_generation
      }

      current_generation = next_generation

    // reset the game if there are no clients
    } else {
      current_generation = game.GosperGliderGun()
    }
  }
}
