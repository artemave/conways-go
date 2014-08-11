package synchronized_broadcaster

import (
	"errors"
	"fmt"

	"code.google.com/p/go-uuid/uuid"
)

type SynchronizedBroadcasterClient interface {
	ClientId() string
	Inbox() chan BroadcastMessage
}

type SynchronizedBroadcaster struct {
	clients      []SynchronizedBroadcasterClient
	messageQueue chan BroadcastMessage
	messageAck   chan bool
}

func NewSynchronizedBroadcaster() SynchronizedBroadcaster {
	sb := SynchronizedBroadcaster{
		clients:      []SynchronizedBroadcasterClient{},
		messageQueue: make(chan BroadcastMessage),
		messageAck:   make(chan bool, 10),
	}

	go func() {
		for {
			select {
			case <-sb.messageAck: // remove possible ack from remove client while sb was idle
				fmt.Printf("Rm client ack\n")
			default:
				for msg := range sb.messageQueue {
					for _, c := range sb.clients {
						c := c
						go func() { c.Inbox() <- msg }()
					}

					ackNum := 0
					for _ = range sb.messageAck {
						fmt.Printf("Ack\n")
						ackNum += 1
						if ackNum >= len(sb.clients) {
							break
						}
					}
				}
				break
			}
		}
	}()

	return sb
}

func (sb SynchronizedBroadcaster) AddClient(client SynchronizedBroadcasterClient) {
	sb.clients = append(sb.clients, client)
}

func (sb SynchronizedBroadcaster) Clients() []SynchronizedBroadcasterClient {
	return sb.clients
}

func (sb SynchronizedBroadcaster) RemoveClient(client SynchronizedBroadcasterClient) error {

	for i, c := range sb.clients {
		if c.ClientId() == client.ClientId() {
			sb.clients = append(sb.clients[:i], sb.clients[i+1:]...)
			sb.messageAck <- true
			return nil
		}
	}

	return errors.New("Trying to remove non existent client")
}

func (sb SynchronizedBroadcaster) MessageAcknowledged() {
	sb.messageAck <- true
}

func (sb SynchronizedBroadcaster) SendBroadcastMessage(data interface{}) {
	u4 := uuid.New()
	msg := BroadcastMessage{
		MessageId: u4,
		Server:    &sb,
		Data:      data,
	}

	sb.messageQueue <- msg
}

type BroadcastMessage struct {
	MessageId string
	Data      interface{}
	Server    *SynchronizedBroadcaster
}
