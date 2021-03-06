package synchronized_broadcaster_test

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
	. "github.com/artemave/conways-go/synchronized_broadcaster"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestClient struct {
	id       string
	Messages chan BroadcastMessage
}

func NewTestClient() *TestClient {
	id := uuid.New()

	tc := &TestClient{
		id:       id,
		Messages: make(chan BroadcastMessage, 1),
	}
	return tc
}

func (tc TestClient) ClientId() string {
	return tc.id
}

func (tc TestClient) Inbox() chan BroadcastMessage {
	return tc.Messages
}

var _ = Describe("SynchronizedBroadcaster", func() {
	var sb *SynchronizedBroadcaster
	var client1 *TestClient
	var client2 *TestClient

	BeforeEach(func() {
		sb = NewSynchronizedBroadcaster()
		client1 = NewTestClient()
		client2 = NewTestClient()

		sb.AddClient(client1)
		sb.AddClient(client2)
	})

	AfterEach(func() {
		for _, c := range sb.Clients() {
			if c.ClientId() == client1.ClientId() {
				sb.RemoveClient(client1)
			}
			if c.ClientId() == client2.ClientId() {
				sb.RemoveClient(client2)
			}
		}
		close(client1.Inbox())
		close(client2.Inbox())
	})

	It("Broadcasts message to clients", func() {
		sb.SendBroadcastMessage("msg")

		msg1 := <-client1.Messages
		Expect(msg1.Data).To(Equal("msg"))
		msg2 := <-client2.Messages
		Expect(msg2.Data).To(Equal("msg"))
	})

	Context("Message in progress", func() {
		assertBlocksUntilAllClientsAcknowledgedMessage := func() {
			d := make(chan string)

			sb.SendBroadcastMessage("msg1")
			sb.MessageAcknowledged(client1)

			// send second message before first one is acknowledged by both clients
			go func() {
				sb.SendBroadcastMessage("msg3")
				d <- (<-client1.Messages).Data.(string)
				d <- (<-client2.Messages).Data.(string)
				d <- (<-client1.Messages).Data.(string)
				d <- (<-client2.Messages).Data.(string)
			}()

			time.Sleep(20 * time.Millisecond)

			go func() {
				d <- "msg2"
			}()

			time.Sleep(20 * time.Millisecond)
			// complete acknowledge first message
			sb.MessageAcknowledged(client2)

			res := []string{<-d, <-d, <-d, <-d, <-d}

			// the order is important
			Expect(res).To(Equal([]string{"msg2", "msg1", "msg1", "msg3", "msg3"}))
		}

		It("Blocks until all clients acknowledged it", assertBlocksUntilAllClientsAcknowledgedMessage)

		Describe("Client disconnects", func() {
			It("Acknowledges message for that client (to prevent from waiting forever)", func(done Done) {
				d := make(chan string)

				sb.SendBroadcastMessage("msg")
				sb.MessageAcknowledged(client2)

				// send second message before first one is acknowledged by both clients
				go func() {
					sb.SendBroadcastMessage("msg1")
					d <- "msg1"
				}()

				sb.RemoveClient(client1)

				Expect(<-d).To(Equal("msg1"))
				close(done)
			}, 0.5)
		})

		Context("Client disconnected while no message was in progress", func() {
			var client3 *TestClient
			BeforeEach(func() {
				client3 = NewTestClient()
				sb.AddClient(client3)
				sb.RemoveClient(client3)
			})

			Context("Another message received", func() {
				It("Still blocks until all clients acknowledged it", assertBlocksUntilAllClientsAcknowledgedMessage)
			})
		})
	})
})
