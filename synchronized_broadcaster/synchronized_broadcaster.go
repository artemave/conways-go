package synchronized_broadcaster

import "code.google.com/p/go-uuid/uuid"

import "github.com/araddon/gou"
import "log"
import "os"

func init() {
	gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")
}

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
				for i, c := range sb.clients {
					if c.ClientId() == client.ClientId() {
						sb.clients = append(sb.clients[:i], sb.clients[i+1:]...)
						sb.clientRemoved <- true
					}
				}
			case msg := <-sb.messageQueue:
				if len(sb.clients) > 0 {

					clientAcks := make(map[string]bool)

					for _, c := range sb.clients {
						clientAcks[c.ClientId()] = true
						c := c
						go func() {
							// don't fail if Inbox is closed
							defer func() { recover() }()
							gou.Debug("Sending message to client ", c.ClientId())
							c.Inbox() <- msg
							gou.Debug("Message sent to client ", c.ClientId())
						}()
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
							if len(clientAcks) == 0 {
								break ACKS
							}
						case client := <-sb.messageAck:
							delete(clientAcks, client.ClientId())
							if len(clientAcks) == 0 {
								break ACKS
							}
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
