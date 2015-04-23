package main

import (
	"net/http"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	. "github.com/artemave/conways-go/game"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

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

	game := gamesRepo.FindGameById(id)
	if game == nil {
		ws.WriteJSON(map[string]string{"Handshake": "game_not_found"})
		return
	}

	player, err := game.AddPlayer()
	if err != nil {
		ws.WriteJSON(map[string]string{"Handshake": "game_taken"})
		return
	} else {
		defer game.RemovePlayer(player)
	}

	if game.IsPractice() {
		practiceWall := NewPracticeWall(game)
		defer practiceWall.RemoveDummyPlayer()
	}

	go Respond(ws, game, player)
	Listen(ws, game, player)
}

type WsServerMessage struct {
	Handshake      string
	Player         int
	Cols           int
	Rows           int
	WinSpots       []WinSpot
	PausedByPlayer int
	FreeCellsCount int
}

type WsServerGameDataMessage struct {
	Handshake      string
	Generation     *conway.Generation
	FreeCellsCount int
}

func Respond(ws *websocket.Conn, game *Game, player *Player) {
	for {
		msg, ok := <-player.GameServerMessages

		if !ok {
			return
		}

		switch messageData := msg.Data.(type) {
		case PlayersAreReady:
			if messageData {
				var serverMessage WsServerMessage

				if game.IsPaused() {
					serverMessage = WsServerMessage{
						Handshake:      "pause",
						Player:         int(player.PlayerIndex),
						PausedByPlayer: int(game.PausedByPlayer),
					}
				} else {
					serverMessage = WsServerMessage{
						Handshake:      "ready",
						Player:         int(player.PlayerIndex),
						FreeCellsCount: int(game.FreeCellsCountOf(player)),
						Cols:           game.Cols(),
						Rows:           game.Rows(),
						WinSpots:       game.WinSpots(),
					}
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
		case PauseGame:
			var serverMessage WsServerMessage
			if messageData {
				serverMessage = WsServerMessage{
					Handshake:      "pause",
					Player:         int(player.PlayerIndex),
					PausedByPlayer: int(game.PausedByPlayer),
				}
			} else {
				serverMessage = WsServerMessage{
					Handshake:      "resume",
					Player:         int(player.PlayerIndex),
					FreeCellsCount: int(game.FreeCellsCountOf(player)),
					Cols:           game.Cols(),
					Rows:           game.Rows(),
					WinSpots:       game.WinSpots(),
				}
			}
			if err := ws.WriteJSON(serverMessage); err != nil {
				gou.Error("Send to user: ", err)
				return
			}
		case *conway.Generation:
			serverMessage := WsServerGameDataMessage{
				Handshake:      "game_data",
				Generation:     messageData,
				FreeCellsCount: int(game.FreeCellsCountOf(player)),
			}
			if err := ws.WriteJSON(serverMessage); err != nil {
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
	}
}

type WsClientMessage struct {
	Acknowledged string        `json:acknowledged,omitempty`
	Command      string        `json:command,omitempty`
	NewCells     []conway.Cell `json:cells,omitempty`
}

func Listen(ws *websocket.Conn, game *Game, player *Player) {
	for {
		var msg WsClientMessage
		if err := ws.ReadJSON(&msg); err != nil {
			return
		} else {
			if msg.Command != "" {
				switch msg.Command {
				case "pause":
					go func() { game.PauseBy(player) }()
				case "resume":
					go func() { game.Resume() }()
				default:
					gou.Errorf("Unknown command %s\n", msg.Command)
				}
			} else if msg.NewCells != nil {
				game.AddCells(msg.NewCells)
			} else {
				switch msg.Acknowledged {
				case "ready", "wait", "game", "finish", "pause", "resume":
					player.MessageAcknowledged()
				default:
					gou.Error("Unknown client message")
				}
			}
		}
	}
}
