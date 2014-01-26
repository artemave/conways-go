package routes

import (
	"code.google.com/p/go.net/websocket"
	"github.com/araddon/gou"
	"github.com/artemave/conways-go/comm"
	"net/http"
)

func InitRoutes(server *comm.Server) {
	/* http.Handle("/", NewGameHandler) */
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
}
