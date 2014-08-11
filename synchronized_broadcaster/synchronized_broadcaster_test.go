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
	var sb SynchronizedBroadcaster
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
		select {
		case <-client1.Inbox():
		case <-client2.Inbox():
		default:
			close(client1.Inbox())
			close(client2.Inbox())
			return
		}
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

			sb.SendBroadcastMessage("msg")
			sb.MessageAcknowledged()

			// send second message before first one is acknowledged by both clients
			go func() {
				sb.SendBroadcastMessage("msg2")
				d <- "msg2"
			}()

			<-time.After(10 * time.Millisecond)

			go func() {
				d <- "msg1"
			}()

			<-time.After(10 * time.Millisecond)
			// complete acknowledge first message
			sb.MessageAcknowledged()

			res := []string{<-d, <-d}

			// the order is important
			Expect(res).To(Equal([]string{"msg1", "msg2"}))
		}

		It("Blocks until all clients acknowledged it", assertBlocksUntilAllClientsAcknowledgedMessage)

		Describe("Client disconnects", func() {
			It("Acknowledges message for that client (to prevent from waiting forever)", func(done Done) {
				d := make(chan string)

				sb.SendBroadcastMessage("msg")
				sb.MessageAcknowledged()

				// send second message before first one is acknowledged by both clients
				go func() {
					sb.SendBroadcastMessage("msg1")
					d <- "msg1"
				}()

				sb.RemoveClient(client1)

				Expect(<-d).To(Equal("msg1"))
				close(done)
			}, 0.05)
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
