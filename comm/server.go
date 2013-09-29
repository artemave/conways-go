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
}

func NewServer() *Server {
  s := &Server{
    addClient: make(chan *Client),
    delClient: make(chan *Client),
    clients: make(map[int]*Client),
    game: game.Game{Rows: 3000, Cols: 4000},
  }
  go s.WatchClientRegister()
  return s
}

func (this *Server) WatchClientRegister() {
  for {
    select {
    case client := <-this.addClient:
      this.clients[client.id] = client
    case client := <-this.delClient:
      delete(this.clients, client.id)
    }
  }
}

func (this *Server) ServeTheGame() {
  current_generation := game.GosperGliderGun()
  for {
    time.Sleep(time.Millisecond * 200)

    if len(this.clients) > 0 {
      next_generation := this.game.NextGeneration(current_generation)

      for _, client := range this.clients {
        client.outToWebClient <- next_generation
      }

      current_generation = next_generation

    // reset the game if there are no clients
    } else {
      current_generation = game.GosperGliderGun()
    }
  }
}
