package game

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/artemave/conways-go/conway"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
)

type Player struct {
	GameServerMessages chan sb.BroadcastMessage
	id                 string
	Broadcaster
	PlayerIndex conway.Player
}

func NewPlayer(b Broadcaster, pi conway.Player) *Player {
	u4 := uuid.New()
	player := &Player{
		Broadcaster:        b,
		PlayerIndex:        pi,
		id:                 u4,
		GameServerMessages: make(chan sb.BroadcastMessage),
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
	p.Broadcaster.MessageAcknowledged(p)
}

func (p Player) CleanUp() {
	close(p.GameServerMessages)
}
