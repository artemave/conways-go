package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/artemave/conways-go/conway"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
)

type Player struct {
	GameServerMessages      chan sb.BroadcastMessage
	id                      string
	SynchronizedBroadcaster *sb.SynchronizedBroadcaster
	PlayerIndex             conway.Player
}

func NewPlayer(game *Game) *Player {
	u4 := uuid.New()
	player := &Player{
		id: u4,
		SynchronizedBroadcaster: game.SynchronizedBroadcaster,
		GameServerMessages:      make(chan sb.BroadcastMessage),
	}
	return player
}

func (p Player) ClientId() string {
	return p.id
}

func (p Player) Inbox() chan sb.BroadcastMessage {
	return p.GameServerMessages
}

func (p Player) MessageAcknowledged() {
	p.SynchronizedBroadcaster.MessageAcknowledged()
}
