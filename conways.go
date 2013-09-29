package main

import (
  "github.com/artemave/conways-go/comm"
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

  server := comm.NewServer()

  go server.ServeTheGame()

  http.Handle("/go-ws", websocket.Handler(func(ws *websocket.Conn) {
    client := comm.NewClient(ws, server)
    gou.Debug("Connected: ", client.Id())

    client.ListenAndServeBackToWebClient()

    defer func() {
      ws.Close()
      gou.Debug("Diconnected: ", client.Id())
    }()
  }))

  http.Handle("/", http.FileServer(http.Dir("./public/")))

  gou.Debug("listening at " + port)
  http.ListenAndServe(":" + port, nil)
}
