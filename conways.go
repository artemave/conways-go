package main

import (
  "github.com/artemave/conways-go/game"
  "github.com/araddon/gou"
  "net/http"
  "os"
  "log"
  "code.google.com/p/go.net/websocket"
  "time"
)


var g = game.Game{Rows: 3000, Cols: 4000}
var selection = make(chan []game.Point)
var current_generation = game.GosperGliderGun()
var errors = make(chan error)

func main() {
  gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

  port := os.Getenv("PORT")
  if (port == "") {
    port = "9999"
  }

  http.Handle("/go-ws", websocket.Handler(GameServerWs))
  http.Handle("/", http.FileServer(http.Dir("./public/")))

  gou.Debug("listening at " + port)
  http.ListenAndServe(":" + port, nil)
}

func GameServerWs(ws *websocket.Conn) {
  defer ws.Close()

  gou.Debug("WS /go-ws")

  go Listen(ws, selection, errors)

  for {
    time.Sleep(time.Millisecond * 100)

    select {
    case points := <- selection:
      current_generation.AddPoints(points)
    case err := <- errors:
      gou.Error(err)
      return
    default:
    }

    next_generation := g.NextGeneration(current_generation)
    websocket.JSON.Send(ws, g.GenerationToPoints(next_generation))

    current_generation = next_generation
  }
}

func Listen(ws *websocket.Conn, selection chan []game.Point, errors chan error) {
  for {
    var points []game.Point
    if err := websocket.JSON.Receive(ws, &points); err != nil {
      errors <- err
      break
    } else {
      for _, point := range points {
        gou.Debug("point from client: ", point, "\n")
      }
      selection <- points
    }
  }
}
