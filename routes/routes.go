package routes

import (
	// "github.com/araddon/gou"
	// "github.com/artemave/conways-go/comm"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes() {
	// server := comm.NewServer()

	// go server.ServeTheGame()

	// http.Handle("/go-ws", websocket.Handler(func(ws *websocket.Conn) {
	// 	client := comm.NewClient(ws, server)
	// 	gou.Debug("Connected: ", client.Id())

	// 	client.ListenAndServeBackToWebClient()

	// 	defer func() {
	// 		ws.Close()
	// 		gou.Debug("Diconnected: ", client.Id())
	// 	}()
	// }))
	r := mux.NewRouter()

	r.HandleFunc("/", StartNewGameHandler)
	r.HandleFunc("/games/{id}", NewGameHandler)
	r.HandleFunc("/games/play/{id}", GamePlayHandler)

	http.Handle("/", r)
	http.Handle("/public/", http.FileServer(http.Dir("./")))
}
