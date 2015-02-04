package main_test

import (
	"fmt"
	"net/http/httptest"
	"time"

	. "github.com/artemave/conways-go"
	"github.com/artemave/conways-go/conway"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var server = httptest.NewServer(nil)

func wsRequest(path string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(httpToWs(server.URL+path), nil)
	if err != nil {
		panic("Dial() returned error " + err.Error())
	}
	return ws
}

var _ = Describe("GamePlayHandler", func() {

	var clockStep int = 10
	*TestDelay = time.Duration(clockStep)

	var startGeneration = &conway.Generation{
		{Point: conway.Point{Row: 3, Col: 2}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 13, Col: 14}, State: conway.Live, Player: conway.Player2},
	}

	BeforeEach(func() {
		(*TestGameRepo).Empty()
		(*TestGameRepo).CreateGameById("123", "small", startGeneration)
	})

	Context("New game", func() {
		It("tells web client to wait", func() {
			ws := wsRequest("/games/play/123")
			defer ws.Close()

			output := justReadHandshake(ws)
			Expect(output.Handshake).To(Equal("wait"))
		})
	})

	Context("Existing game", func() {
		var firstWs *websocket.Conn
		var secondWs *websocket.Conn

		BeforeEach(func() {
			firstWs = wsRequest("/games/play/123")
			justReadHandshake(firstWs)
			sendAckMessage(firstWs, "wait")
		})
		AfterEach(func() {
			firstWs.Close()
		})

		Context("game is paused", func() {
			BeforeEach(func() {
				sendPause(firstWs)
				justReadHandshake(firstWs)
				sendAckMessage(firstWs, "pause")
			})
			Describe("second client connects", func() {
				BeforeEach(func() { secondWs = wsRequest("/games/play/123") })
				AfterEach(func() { secondWs.Close() })

				It("tells second client that the game is paused", func() {
					output := justReadHandshake(secondWs)
					Expect(output.Handshake).To(Equal("pause"))
				})

				It("tells web client what player paused the game", func() {
					output := justReadHandshake(secondWs)
					Expect(output.PausedByPlayer).To(Equal(1))
				})
			})
		})

		Describe("second client connects", func() {
			BeforeEach(func() { secondWs = wsRequest("/games/play/123") })
			AfterEach(func() { secondWs.Close() })

			It("tells all web clients to join the game", func() {
				output := justReadHandshake(secondWs)
				Expect(output.Handshake).To(Equal("ready"))
				output = justReadHandshake(firstWs)
				Expect(output.Handshake).To(Equal("ready"))
			})

			It("tells all web clients their player number", func() {
				output := justReadHandshake(firstWs)
				Expect(output.Player).To(Equal(1))
				output = justReadHandshake(secondWs)
				Expect(output.Player).To(Equal(2))
			})

			It("tells all web clients the field size", func() {
				output := justReadHandshake(firstWs)
				Expect(output.Cols).To(Equal(40))
				Expect(output.Rows).To(Equal(26))
				output = justReadHandshake(secondWs)
				Expect(output.Cols).To(Equal(40))
				Expect(output.Rows).To(Equal(26))
			})

			Context("all clients acknowledged ready", func() {

				BeforeEach(func() {
					justReadHandshake(firstWs)
					justReadHandshake(secondWs)
					sendAckMessage(firstWs, "ready")
					sendAckMessage(secondWs, "ready")
				})

				It("starts serving game to all clients", func() {
					assertGenerationOutput(firstWs)
					assertGenerationOutput(secondWs)
				})

				Describe("second client disconnects", func() {

					var recoverFromExtraGameMessage = func() {
						if r := recover(); r != nil {
							// Recovering from reading wrong type of message
							sendAckMessage(firstWs, "game")
							// Why wrong type? Because if client disconnects at exactly the same time as
							// as game message sent we might get an extra one
						}
					}

					BeforeEach(func() {
						secondWs.Close()
						justReadGameOutput(firstWs)
						sendAckMessage(firstWs, "game")
					})

					It("tells first client to wait", func() {
						Eventually(func() string {
							defer recoverFromExtraGameMessage()
							o := justReadHandshake(firstWs)
							return o.Handshake
						}).Should(Equal("wait"))
					})

					It("stops game broadcast", func() {
						Eventually(func() bool {
							defer recoverFromExtraGameMessage()
							justReadHandshake(firstWs)
							return true
						}).Should(Equal(true))

						sendAckMessage(firstWs, "wait")

						msgSent := make(chan *conway.Generation)
						go func(c chan *conway.Generation) {
							defer func() {
								if r := recover(); r != nil {
									// reading closed ws after test finish should not fail the test
								}
							}()
							o := justReadGameOutput(firstWs)
							c <- o
						}(msgSent)

						for {
							select {
							case g := <-msgSent:
								fmt.Printf("%#v\n", *g)
								Fail("Expected to stop broadcasting game")
							case <-time.After(time.Millisecond * time.Duration(clockStep+10)):
								close(msgSent)
								return
							}
						}
					})
				})
			})

			Context("clients acknowledged game message", func() {
				BeforeEach(func() {
					justReadHandshake(firstWs)
					justReadHandshake(secondWs)

					sendAckMessage(firstWs, "ready")
					sendAckMessage(secondWs, "ready")

					justReadGameOutput(firstWs)
					justReadGameOutput(secondWs)

					sendAckMessage(firstWs, "game")
					sendAckMessage(secondWs, "game")
				})

				It("sends next generation to all clients", func() {
					assertGenerationTwo(firstWs)
					assertGenerationTwo(secondWs)
				})
			})

			Context("third client", func() {
				It("tells web client that the game has already started", func() {
					ws := wsRequest("/games/play/123")
					defer ws.Close()

					output := justReadHandshake(ws)
					Expect(output.Handshake).To(Equal("game_taken"))
				})
			})
		})

		Context("player paused the game", func() {
			BeforeEach(func() {
				*TestDelay = time.Duration(clockStep * 5)
				secondWs = wsRequest("/games/play/123")

				justReadHandshake(firstWs)
				justReadHandshake(secondWs)

				sendAckMessage(firstWs, "ready")
				sendAckMessage(secondWs, "ready")

				justReadGameOutput(firstWs)
				justReadGameOutput(secondWs)

				sendPause(secondWs)
				sendAckMessage(firstWs, "game")
				sendAckMessage(secondWs, "game")
			})
			AfterEach(func() {
				*TestDelay = time.Duration(clockStep)
				secondWs.Close()
			})

			It("tells all players that the game is paused", func() {
				output := justReadHandshake(firstWs)
				Expect(output.Handshake).To(Equal("pause"))

				output = justReadHandshake(secondWs)
				Expect(output.Handshake).To(Equal("pause"))
			})

			Context("player resumes the game", func() {
				BeforeEach(func() {
					justReadHandshake(firstWs)
					justReadHandshake(secondWs)

					sendAckMessage(firstWs, "pause")
					sendAckMessage(secondWs, "pause")

					sendResume(firstWs)
				})

				It("tells all player that the game is resumed", func() {
					output := justReadHandshake(firstWs)
					Expect(output.Handshake).To(Equal("resume"))

					output = justReadHandshake(secondWs)
					Expect(output.Handshake).To(Equal("resume"))
				})

				It("includes game metadata in the message (in case we started paused)", func() {
					game := (*TestGameRepo).FindGameById("123")
					output := justReadHandshake(firstWs)

					Expect(output.Player).To(Equal(1))
					Expect(output.Cols).To(Equal(game.Cols()))
					Expect(output.Rows).To(Equal(game.Rows()))
					Expect(output.WinSpots).To(Equal(game.WinSpots()))
				})
			})
		})
	})
})

func justReadHandshake(ws *websocket.Conn) WsServerMessage {
	var output WsServerMessage
	if err := ws.ReadJSON(&output); err != nil {
		panic(err)
	}
	return output
}

func justReadGameOutput(ws *websocket.Conn) *conway.Generation {
	var output WsServerGameDataMessage
	if err := ws.ReadJSON(&output); err != nil {
		panic(err)
	}
	return output.Generation
}

func assertGenerationOutput(ws *websocket.Conn) {
	justReadGameOutput(ws)
}

func assertGenerationTwo(ws *websocket.Conn) {
	output := justReadGameOutput(ws)

	secondGeneration := &conway.Generation{
		{Point: conway.Point{Row: 2, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
	}

	Expect(output).To(Equal(secondGeneration))
}

func sendAckMessage(ws *websocket.Conn, msg string) {
	if err := ws.WriteJSON(map[string]string{"acknowledged": msg}); err != nil {
		Fail(err.Error())
	}
}

func sendPause(ws *websocket.Conn) {
	if err := ws.WriteJSON(map[string]string{"command": "pause"}); err != nil {
		Fail(err.Error())
	}
}

func sendResume(ws *websocket.Conn) {
	if err := ws.WriteJSON(map[string]string{"command": "resume"}); err != nil {
		Fail(err.Error())
	}
}

func httpToWs(u string) string {
	return "ws" + u[len("http"):]
}
