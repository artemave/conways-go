package main

import (
	"errors"
	"net/http"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/conway"
	"github.com/artemave/conways-go/dependencies/gouuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type SynchronizedBroadcasterClient interface {
	Id() *uuid.UUID
	Inbox(*BroadcastMessage)
	MessageAcknowledged(*uuid.UUID)
}

type SynchronizedBroadcaster struct {
	Clients []SynchronizedBroadcasterClient
}

func NewSynchronizedBroadcaster() *SynchronizedBroadcaster {
	return &SynchronizedBroadcaster{Clients: []SynchronizedBroadcasterClient{}}
}

func (sb *SynchronizedBroadcaster) AddClient(client SynchronizedBroadcasterClient) {
	sb.Clients = append(sb.Clients, client)
}

func (sb *SynchronizedBroadcaster) RemoveClient(client SynchronizedBroadcasterClient) error {

	for i, c := range sb.Clients {
		if c.Id() == client.Id() {
			sb.Clients = append(sb.Clients[:i], sb.Clients[i+1:]...)
			return nil
		}
	}

	return errors.New("Trying to remove non existent client")
}

func (sb *SynchronizedBroadcaster) MessageAcknowledged(msgId *uuid.UUID) error {
	return nil
}

func (sb *SynchronizedBroadcaster) NewBroadcastMessage() *BroadcastMessage {
	u4, _ := uuid.NewV4()
	msg := &BroadcastMessage{
		Id:     u4,
		Server: sb,
	}
	// TODO tell server that message is in progress
	return msg
}

type BroadcastMessage struct {
	Id     *uuid.UUID
	Data   interface{}
	Server *SynchronizedBroadcaster
}

func (bm *BroadcastMessage) SetData(data interface{}) {
	bm.Data = data
}

func (bm *BroadcastMessage) Send() error {
	for _, c := range bm.Server.Clients {
		go c.Inbox(bm)
	}
	return nil
}

func (bm *BroadcastMessage) Discard() {
	//TODO reset server's "message in progress"
}

type Player struct {
	GameServerMessages chan *BroadcastMessage
	id                 *uuid.UUID
	Game               *Game
}

func NewPlayer(game *Game) *Player {
	u4, _ := uuid.NewV4()
	player := &Player{
		id:                 u4,
		Game:               game,
		GameServerMessages: make(chan *BroadcastMessage),
	}
	return player
}

func (p *Player) Id() *uuid.UUID {
	return p.id
}

func (p *Player) Inbox(msg *BroadcastMessage) {
	p.GameServerMessages <- msg
}

func (p *Player) MessageAcknowledged(msgId *uuid.UUID) {
	p.Game.MessageAcknowledged(msgId)
}

type Game struct {
	Id                      string
	SynchronizedBroadcaster *SynchronizedBroadcaster
	Conway                  *conway.Game
}

func NewGame(id string) *Game {
	game := &Game{
		Id: id,
		SynchronizedBroadcaster: NewSynchronizedBroadcaster(),
		Conway:                  &conway.Game{Cols: 300, Rows: 200},
	}
	return game
}

func (g *Game) MessageAcknowledged(msgId *uuid.UUID) {
	g.SynchronizedBroadcaster.MessageAcknowledged(msgId)
}

func (g *Game) AddPlayer() (*Player, error) {
	if len(g.SynchronizedBroadcaster.Clients) >= 2 {
		return &Player{}, errors.New("Game has already reached maximum number players")
	}
	p := NewPlayer(g)

	msg := g.SynchronizedBroadcaster.NewBroadcastMessage()

	g.SynchronizedBroadcaster.AddClient(p)
	enoughPlayersToStart := len(g.SynchronizedBroadcaster.Clients) >= 2

	msg.SetData(enoughPlayersToStart)
	msg.Send()

	if enoughPlayersToStart {
		msg = g.SynchronizedBroadcaster.NewBroadcastMessage()
		msg.SetData(g.NextGeneration())
		msg.Send()
	}

	return p, nil
}

func (g *Game) NextGeneration() *conway.Generation {
	return g.StartGeneration()
}

func (g *Game) StartGeneration() *conway.Generation {
	// 150x100
	return &conway.Generation{
		{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},

		{Point: conway.Point{Row: 95, Col: 145}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 96, Col: 145}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 96, Col: 146}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 95, Col: 146}, State: conway.Live, Player: conway.Player2},
	}
}

func (g *Game) RemovePlayer(p *Player) error {
	msg := g.SynchronizedBroadcaster.NewBroadcastMessage()

	if err := g.SynchronizedBroadcaster.RemoveClient(p); err != nil {
		msg.Discard()
		return errors.New("Trying to delete non-existent player")
	}

	msg.SetData(len(g.SynchronizedBroadcaster.Clients) >= 2)
	msg.Send()

	return nil
}

type GamesRepo struct {
	Games []*Game
}

func NewGamesRepo() *GamesRepo {
	gr := &GamesRepo{
		Games: []*Game{},
	}
	return gr
}

// FIXME this is not thread-safe
func (gr *GamesRepo) FindOrCreateGameById(id string) *Game {
	for _, game := range gr.Games {
		if game.Id == id {
			return game
		}
	}
	newGame := NewGame(id)
	gr.Games = append(gr.Games, newGame)
	return newGame
}

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

	go func() {
		for {
			if _, _, err := ws.NextReader(); err != nil {
				disconnected <- true
			}
		}
	}()

	for {
		select {
		case msg := <-player.GameServerMessages:

			switch messageData := msg.Data.(type) {
			case bool:
				if messageData {
					if err := ws.WriteJSON(map[string]string{"handshake": "ready"}); err != nil {
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
				points := conway.GenerationToPoints(messageData)
				if err := ws.WriteJSON(points); err != nil {
					gou.Error("Send to user: ", err)
					return
				}
			}

			player.MessageAcknowledged(msg.Id)
		case <-disconnected:
			return
		}
	}
}
