package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/araddon/gou"
	"github.com/artemave/conways-go/comm"
	"log"
	"net/http"
	"os"
)

func main() {
	gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

	port := os.Getenv("PORT")
	if port == "" {
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

	http.Handle("/public/", http.FileServer(http.Dir("./")))

	gou.Debug("listening at " + port)
	http.ListenAndServe(":"+port, nil)
}
