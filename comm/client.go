package comm

import (
  "github.com/artemave/conways-go/game"
  "code.google.com/p/go.net/websocket"
  "github.com/araddon/gou"
)

type Client struct {
  id int
  server *Server
  ws *websocket.Conn
  outToUser chan *game.Generation
  disconnect chan bool
}

var maxId int = 0

func NewClient(ws *websocket.Conn, server *Server) *Client {
  maxId++

  client := &Client{
    id: maxId,
    server: server,
    ws: ws,
    outToUser: make(chan *game.Generation),
    disconnect: make(chan bool),
  }
  return client
}

func (this *Client) Id() int {
  return this.id
}

func (this *Client) ListenAndServeBackToWebClient() {
  this.server.addClient <- this

  go this.ListenUserEvents()
  go this.WriteBackToUser()

  <-this.disconnect
  this.server.delClient <- this
}

func (this *Client) WriteBackToUser() {
  for {
    g := <-this.outToUser

    points := game.GenerationToPoints(g)

    if err := websocket.JSON.Send(this.ws, points); err != nil {
      gou.Error("Send to user: ", err)
      this.disconnect <- true
      break
    }
  }
}

func (this *Client) ListenUserEvents() {
  for {
    var points []game.Point

    if err := websocket.JSON.Receive(this.ws, &points); err != nil {
      gou.Error("Receive from user: ", err)
      this.disconnect <- true
      break
    }
    this.server.pointsFromUsers <- points
  }
}
