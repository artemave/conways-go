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
	PlayerIndex    conway.Player
	freeCellsCount *CellCount
}

type CellCount int

const maxFreeCells = 10

func NewPlayer(g *Game) *Player {
	initMaxFreeCells := CellCount(maxFreeCells)

	u4 := uuid.New()
	player := &Player{
		id:                 u4,
		Broadcaster:        g.Broadcaster,
		GameServerMessages: make(chan sb.BroadcastMessage),
		freeCellsCount:     &initMaxFreeCells,
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

func (p Player) FreeCellsCount() CellCount {
	return *p.freeCellsCount
}

func (p Player) DecreaseFreeCellsCountBy(count int) {
	*p.freeCellsCount = CellCount(int(*p.freeCellsCount) - count)
}

func (p Player) NextFreeCellsCount() CellCount {
	if *p.freeCellsCount < maxFreeCells {
		*p.freeCellsCount = *p.freeCellsCount + 1
	}
	return *p.freeCellsCount
}
