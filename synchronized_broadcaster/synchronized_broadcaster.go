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
	clients           []SynchronizedBroadcasterClient
	addClient         chan SynchronizedBroadcasterClient
	clientAdded       chan bool
	removeClient      chan SynchronizedBroadcasterClient
	clientRemoved     chan bool
	clientInboxClosed map[string]chan bool
	messageQueue      chan BroadcastMessage
	messageAck        chan SynchronizedBroadcasterClient
}

func NewSynchronizedBroadcaster() *SynchronizedBroadcaster {
	sb := &SynchronizedBroadcaster{
		clients:           []SynchronizedBroadcasterClient{},
		addClient:         make(chan SynchronizedBroadcasterClient),
		clientAdded:       make(chan bool),
		removeClient:      make(chan SynchronizedBroadcasterClient),
		clientRemoved:     make(chan bool),
		clientInboxClosed: make(map[string]chan bool),
		messageQueue:      make(chan BroadcastMessage),
		messageAck:        make(chan SynchronizedBroadcasterClient),
	}

	go func() {
		for {
			select {
			case client := <-sb.addClient:
				sb.clients = append(sb.clients, client)
				sb.clientAdded <- true
			case client := <-sb.removeClient:
				gou.Debug("Rm client (1)", client.ClientId())
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
						gou.Debug("Message to client", c.ClientId())
						go func() {
							select {
							case <-sb.clientInboxClosed[c.ClientId()]:
								gou.Debug("Skipping message to client", c.ClientId())
							default:
								gou.Debug("Sending message to client", c.ClientId())
								c.Inbox() <- msg
							}
						}()
					}

				ACKS:
					for {
						select {
						case client := <-sb.removeClient:
							gou.Debug("Rm client (2)", client.ClientId())
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
	sb.clientInboxClosed[client.ClientId()] = make(chan bool, 1)
	sb.addClient <- client
	<-sb.clientAdded
}

func (sb *SynchronizedBroadcaster) Clients() []SynchronizedBroadcasterClient {
	return sb.clients
}

func (sb *SynchronizedBroadcaster) RemoveClient(client SynchronizedBroadcasterClient) {
	gou.Debug("RemoveClient start", client.ClientId())
	sb.removeClient <- client
	sb.clientInboxClosed[client.ClientId()] <- true
EMPTY_CLIENT_INBOX:
	for {
		select {
		case <-client.Inbox():
			gou.Debug("Draining inbox of client", client.ClientId())
		default:
			gou.Debug("Send client removed", client.ClientId())
			<-sb.clientRemoved
			close(sb.clientInboxClosed[client.ClientId()])
			delete(sb.clientInboxClosed, client.ClientId())
			break EMPTY_CLIENT_INBOX
		}
	}
	gou.Debug("Client removed", client.ClientId())
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
