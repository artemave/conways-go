package synchronized_broadcaster

import (
	"errors"

	"github.com/nu7hatch/gouuid"
)

type SynchronizedBroadcasterClient interface {
	ClientId() *uuid.UUID
	Inbox() chan BroadcastMessage
}

type SynchronizedBroadcaster struct {
	Clients      []SynchronizedBroadcasterClient
	messageQueue chan BroadcastMessage
	messageAck   chan *uuid.UUID
}

func NewSynchronizedBroadcaster() *SynchronizedBroadcaster {
	sb := &SynchronizedBroadcaster{
		Clients:      []SynchronizedBroadcasterClient{},
		messageQueue: make(chan BroadcastMessage),
		messageAck:   make(chan *uuid.UUID),
	}

	go func() {
		for msg := range sb.messageQueue {
			for _, c := range sb.Clients {
				c := c
				go func() {
					c.Inbox() <- msg
				}()
			}

			ackNum := 0
			for _ = range sb.messageAck {
				ackNum += 1
				if ackNum == len(sb.Clients) {
					break
				}
			}
		}
	}()

	return sb
}

func (sb *SynchronizedBroadcaster) AddClient(client SynchronizedBroadcasterClient) {
	sb.Clients = append(sb.Clients, client)
}

func (sb *SynchronizedBroadcaster) RemoveClient(client SynchronizedBroadcasterClient) error {

	for i, c := range sb.Clients {
		if c.ClientId() == client.ClientId() {
			sb.Clients = append(sb.Clients[:i], sb.Clients[i+1:]...)
			close(client.Inbox())
			return nil
		}
	}

	return errors.New("Trying to remove non existent client")
}

func (sb *SynchronizedBroadcaster) MessageAcknowledged(client SynchronizedBroadcasterClient) {
	sb.messageAck <- client.ClientId()
}

func (sb *SynchronizedBroadcaster) SendBroadcastMessage(data interface{}) {
	u4, _ := uuid.NewV4()
	msg := BroadcastMessage{
		MessageId: u4,
		Server:    sb,
		Data:      data,
	}

	sb.messageQueue <- msg
}

type BroadcastMessage struct {
	MessageId *uuid.UUID
	Data      interface{}
	Server    *SynchronizedBroadcaster
}
