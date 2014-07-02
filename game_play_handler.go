package main

import (
	"fmt"
	"net/http"
	"strconv"
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
	{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 4, Col: 3}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},

	{Point: conway.Point{Row: 95, Col: 145}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 96, Col: 145}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 96, Col: 146}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 95, Col: 146}, State: conway.Live, Player: conway.Player2},
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

	go Listen(ws, player, disconnected)
	Respond(ws, game, player, disconnected)
}

func Respond(ws *websocket.Conn, game *Game, player *Player, disconnected chan bool) {
	for {
		select {
		case msg := <-player.GameServerMessages:

			switch messageData := msg.Data.(type) {
			case bool:
				if messageData {
					if err := ws.WriteJSON(map[string]string{"handshake": "ready", "player": strconv.Itoa(game.PlayerNumber(player))}); err != nil {
						gou.Error("Send to user: ", err)
						return
					}
				} else {
					if err := ws.WriteJSON(map[string]string{"handshake": "wait"}); err != nil {
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

func Listen(ws *websocket.Conn, player *Player, disconnected chan bool) {
	for {
		var msg map[string]interface{}
		if err := ws.ReadJSON(&msg); err != nil {
			disconnected <- true
			return
		} else {
			switch msg["acknowledged"].(string) {
			case "ready", "wait", "game":
				player.MessageAcknowledged()
			default:
				fmt.Printf("Unknown client message\n")
			}
		}
	}
}
