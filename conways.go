package main

import (
  "github.com/artemave/conways-go/game"
  "github.com/araddon/gou"
  "net/http"
  "os"
  "log"
  "code.google.com/p/go.net/websocket"
)

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

  var g = game.Game{Rows: 3000, Cols: 4000}

  for {
    var points []game.Point
    if err := websocket.JSON.Receive(ws, &points); err != nil {
      gou.Error(err)
      return
    } else {
      current_generation := g.PointsToGeneration(&points)
      next_generation := g.NextGeneration(current_generation)

      websocket.JSON.Send(ws, g.GenerationToPoints(next_generation))
    }
  }
}
