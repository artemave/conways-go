package main_test

import (
	. "github.com/artemave/conways-go"
	sb "github.com/artemave/conways-go/synchronized_broadcaster"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type StubBroadcaster struct{}

func (b StubBroadcaster) Clients() []sb.SynchronizedBroadcasterClient {
	return []sb.SynchronizedBroadcasterClient{}
}
func (b StubBroadcaster) AddClient(sb.SynchronizedBroadcasterClient) {}
func (b StubBroadcaster) RemoveClient(sb.SynchronizedBroadcasterClient) error {
	return nil
}
func (b StubBroadcaster) SendBroadcastMessage(data interface{}) {}
func (b StubBroadcaster) MessageAcknowledged()                  {}

var _ = Describe("Game", func() {
	Describe("CalculateGameResult", func() {
		Context("Game has not started yet", func() {
			It("returns nil", func() {
				game := NewGame("id", "small", (*TestStartGeneration)["small"])
				game.Broadcaster = StubBroadcaster{}
				game.AddPlayer()
				game.AddPlayer()

				res := game.CalculateGameResult()
				Expect(res).To(BeNil())
			})
		})
	})
})
