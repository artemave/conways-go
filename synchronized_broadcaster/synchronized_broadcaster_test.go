package synchronized_broadcaster_test

import (
	"time"
	. "github.com/artemave/conways-go/synchronized_broadcaster"
	"github.com/nu7hatch/gouuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestClient struct {
	id       *uuid.UUID
	Messages chan BroadcastMessage
}

func NewTestClient() *TestClient {
	id, _ := uuid.NewV4()

	tc := &TestClient{
		id:       id,
		Messages: make(chan BroadcastMessage),
	}
	return tc
}

func (tc TestClient) ClientId() *uuid.UUID {
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

	It("Broadcasts message to clients", func() {
		sb.SendBroadcastMessage("msg")

		msg1 := <-client1.Messages
		Expect(msg1.Data).To(Equal("msg"))
		msg2 := <-client2.Messages
		Expect(msg2.Data).To(Equal("msg"))
	})

	Context("Message in progress", func() {
		It("Blocks until all clients acknowledged it", func() {
			d := make(chan string)

			sb.SendBroadcastMessage("msg")
			sb.MessageAcknowledged(client2)

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
			sb.MessageAcknowledged(client1)

			res := []string{<-d, <-d}

			// the order is important
			Expect(res).To(Equal([]string{"msg1", "msg2"}))
		})
	})
})
