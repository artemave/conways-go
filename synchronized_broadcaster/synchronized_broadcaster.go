package synchronized_broadcaster

import (
	"fmt"

	"code.google.com/p/go-uuid/uuid"
)

type SynchronizedBroadcasterClient interface {
	ClientId() string
	Inbox() chan BroadcastMessage
}

type SynchronizedBroadcaster struct {
	clients       []SynchronizedBroadcasterClient
	addClient     chan SynchronizedBroadcasterClient
	clientAdded   chan bool
	removeClient  chan SynchronizedBroadcasterClient
	clientRemoved chan bool
	messageQueue  chan BroadcastMessage
	messageAck    chan SynchronizedBroadcasterClient
}

func NewSynchronizedBroadcaster() *SynchronizedBroadcaster {
	sb := &SynchronizedBroadcaster{
		clients:       []SynchronizedBroadcasterClient{},
		addClient:     make(chan SynchronizedBroadcasterClient),
		clientAdded:   make(chan bool),
		removeClient:  make(chan SynchronizedBroadcasterClient),
		clientRemoved: make(chan bool),
		messageQueue:  make(chan BroadcastMessage),
		messageAck:    make(chan SynchronizedBroadcasterClient),
	}

	go func() {
		for {
			select {
			case client := <-sb.addClient:
				sb.clients = append(sb.clients, client)
				sb.clientAdded <- true
			case client := <-sb.removeClient:
				fmt.Printf("FUCK\n")
				for i, c := range sb.clients {
					if c.ClientId() == client.ClientId() {
						sb.clients = append(sb.clients[:i], sb.clients[i+1:]...)
						sb.clientRemoved <- true
					}
				}
			case msg := <-sb.messageQueue:
				clientAcks := make(map[string]bool)

				for _, c := range sb.clients {
					clientAcks[c.ClientId()] = true
					c := c
					go func() { c.Inbox() <- msg }()
				}

			ACKS:
				for {
					select {
					case client := <-sb.removeClient:
						for i, c := range sb.clients {
							if c.ClientId() == client.ClientId() {
								sb.clients = append(sb.clients[:i], sb.clients[i+1:]...)
								sb.clientRemoved <- true
								delete(clientAcks, c.ClientId())
							}
						}
					case client := <-sb.messageAck:
						delete(clientAcks, client.ClientId())
					default:
						if len(clientAcks) == 0 {
							break ACKS
						}
					}
				}
			}
		}
	}()

	return sb
}

func (sb *SynchronizedBroadcaster) AddClient(client SynchronizedBroadcasterClient) {
	sb.addClient <- client
	<-sb.clientAdded
}

func (sb *SynchronizedBroadcaster) Clients() []SynchronizedBroadcasterClient {
	return sb.clients
}

func (sb *SynchronizedBroadcaster) RemoveClient(client SynchronizedBroadcasterClient) {
	sb.removeClient <- client
	<-sb.clientRemoved
}

func (sb *SynchronizedBroadcaster) MessageAcknowledged(client SynchronizedBroadcasterClient) {
	sb.messageAck <- client
}

func (sb *SynchronizedBroadcaster) SendBroadcastMessage(data interface{}) {
	u4 := uuid.New()
	msg := BroadcastMessage{
		MessageId: u4,
		Server:    sb,
		Data:      data,
	}

	sb.messageQueue <- msg
}

type BroadcastMessage struct {
	MessageId string
	Data      interface{}
	Server    *SynchronizedBroadcaster
}
