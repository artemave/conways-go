package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var delay = time.Duration(1000)

var gamesRepo = NewGamesRepo()

// 150x100
var startGeneration = &conway.Generation{
	{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 5, Col: 4}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 5, Col: 5}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 4, Col: 5}, State: conway.Live, Player: conway.Player1},

	{Point: conway.Point{Row: 44, Col: 73}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 45, Col: 73}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 45, Col: 74}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 44, Col: 74}, State: conway.Live, Player: conway.Player2},
}

func GamePlayHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		gou.Error(err)
		return
	}
	defer ws.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	gou.Debug("/games/%v", id)

	game := gamesRepo.FindOrCreateGameById(id)

	player, err := game.AddPlayer()

	if err != nil {
		ws.WriteJSON(map[string]string{"handshake": "game_taken"})
		return
	} else {
		defer game.RemovePlayer(player)
	}

	disconnected := make(chan bool)

	go Listen(ws, game, player, disconnected)
	Respond(ws, game, player, disconnected)
}

type WsServerMessage struct {
	Handshake string `json:handshake`
	Player    int    `json:player`
}

func Respond(ws *websocket.Conn, game *Game, player *Player, disconnected chan bool) {
	for {
		select {
		case msg := <-player.GameServerMessages:

			switch messageData := msg.Data.(type) {
			case bool:
				if messageData {
					if err := ws.WriteJSON(WsServerMessage{Handshake: "ready", Player: game.PlayerNumber(player)}); err != nil {
						gou.Error("Send to user: ", err)
						return
					}
				} else {
					if err := ws.WriteJSON(WsServerMessage{Handshake: "wait"}); err != nil {
						gou.Error("Send to user: ", err)
						return
					}
				}
			case *conway.Generation:
				if err := ws.WriteJSON(messageData); err != nil {
					gou.Error("Send to user: ", err)
					return
				}
			}
		case <-disconnected:
			fmt.Printf("Client disconnected\n")
			ws.Close()
			return
		}
	}
}

type WsClientMessage struct {
	Acknowledged string        `json:acknowledged`
	Cells        []conway.Cell `json:cells,omitempty`
}

func Listen(ws *websocket.Conn, game *Game, player *Player, disconnected chan bool) {
	for {
		var msg WsClientMessage
		if err := ws.ReadJSON(&msg); err != nil {
			disconnected <- true
			return
		} else {
			switch msg.Acknowledged {
			case "ready", "wait", "game":
				if msg.Cells != nil {
					// TODO test
					game.AddCells(msg.Cells)
				}
				player.MessageAcknowledged()
			default:
				fmt.Printf("Unknown client message\n")
			}
		}
	}
}
