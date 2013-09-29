package comm

import (
  "github.com/artemave/conways-go/game"
  "code.google.com/p/go.net/websocket"
)

type Client struct {
  id int
  server *Server
  ws *websocket.Conn
  outToWebClient chan *game.Generation
}

var maxId int = 0

func NewClient(ws *websocket.Conn, server *Server) *Client {
  maxId++

  client := &Client{
    id: maxId,
    server: server,
    ws: ws,
    outToWebClient: make(chan *game.Generation),
  }
  return client
}

func (this *Client) Id() int {
  return this.id
}

func (this *Client) ListenAndServeBackToWebClient() {
  this.server.addClient <- this

  for {
    g := <-this.outToWebClient
    points := game.GenerationToPoints(g)
    if err := websocket.JSON.Send(this.ws, points); err != nil {
      this.server.delClient <- this
      break
    }
  }
}
