package main

import (
	"fmt"
	"net/http"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	. "github.com/artemave/conways-go/game"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var gamesRepo = NewGamesRepo()

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

	gou.Debug("WS: /games/play/" + id)

	game := gamesRepo.FindGameById(id)
	if game == nil {
		ws.WriteJSON(map[string]string{"handshake": "game not found"})
		return
	}

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
	Handshake string
	Player    int
	Cols      int
	Rows      int
	WinSpot   conway.Point
}

func Respond(ws *websocket.Conn, game *Game, player *Player, disconnected chan bool) {
	for {
		select {
		case msg := <-player.GameServerMessages:

			switch messageData := msg.Data.(type) {
			case bool:
				if messageData {
					serverMessage := WsServerMessage{
						Handshake: "ready",
						Player:    int(player.PlayerIndex),
						Cols:      game.Cols(),
						Rows:      game.Rows(),
						WinSpot:   *game.WinSpot(&player.PlayerIndex),
					}
					if err := ws.WriteJSON(serverMessage); err != nil {
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
			case GameResult:
				var result string
				switch messageData.Winner {
				case player:
					result = "won"
				case &Player{}:
					result = "draw"
				default:
					result = "lost"
				}
				if err := ws.WriteJSON(map[string]string{"Result": result, "Handshake": "finish"}); err != nil {
					gou.Error("Send to user: ", err)
					return
				}
			}
		case <-disconnected:
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
			case "ready", "wait", "game", "finish":
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
