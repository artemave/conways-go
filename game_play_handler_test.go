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

	var clockStep int = 200
	*TestDelay = time.Duration(clockStep)

	var startGeneration = &conway.Generation{
		{Point: conway.Point{Row: 3, Col: 2}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},
	}

	BeforeEach(func() {
		TestGameRepo.Empty()
		TestGameRepo.CreateGameById("123", "small", startGeneration)
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

		Context("second client connects", func() {
			BeforeEach(func() {
				secondWs = wsRequest("/games/play/123")
			})
			AfterEach(func() {
				secondWs.Close()
			})

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

					BeforeEach(func() {
						// to prevent sending ack to closed channel
						time.Sleep(time.Millisecond * time.Duration(clockStep-20))
						secondWs.Close()
					})

					It("tells first client to wait", func() {
						justReadGameOutput(firstWs)

						output := justReadHandshake(firstWs)
						Expect(output.Handshake).To(Equal("wait"))
					})

					It("stops game broadcast", func() {
						msgSent := make(chan bool)

						justReadGameOutput(firstWs)
						justReadHandshake(firstWs)

						go func(c chan bool) {
							defer func() {
								if r := recover(); r != nil {
									fmt.Println("Recovered after reading closed ws: ", r)
								}
							}()

							justReadGameOutput(firstWs)
							c <- true
						}(msgSent)

						sendAckMessage(firstWs, "wait")

						for {
							select {
							case <-msgSent:
								Fail("Expected to stop broadcasting game")
							case <-time.After(time.Millisecond * time.Duration(clockStep+10)):
								close(msgSent)
								return
							}
						}
					})
				})
			})
		})

		Context("clients acknowledged game message", func() {
			BeforeEach(func() {
				secondWs = wsRequest("/games/play/123")

				justReadHandshake(firstWs)
				justReadHandshake(secondWs)

				sendAckMessage(firstWs, "ready")
				sendAckMessage(secondWs, "ready")

				justReadGameOutput(firstWs)
				justReadGameOutput(secondWs)

				sendAckMessage(firstWs, "game")
				sendAckMessage(secondWs, "game")
			})

			AfterEach(func() {
				secondWs.Close()
			})

			It("sends next generation to all clients", func() {
				assertGenerationTwo(firstWs)
				assertGenerationTwo(secondWs)
			})
		})

		Context("third client", func() {
			var secondWs *websocket.Conn

			BeforeEach(func() {
				secondWs = wsRequest("/games/play/123")
			})
			AfterEach(func() {
				secondWs.Close()
			})

			It("tells web client that the game has already started", func() {
				ws := wsRequest("/games/play/123")
				defer ws.Close()

				output := justReadHandshake(ws)
				Expect(output.Handshake).To(Equal("game_taken"))
			})
		})

	})
})

func justReadHandshake(ws *websocket.Conn) WsServerMessage {
	var output WsServerMessage
	err := ws.ReadJSON(&output)
	if err != nil {
		panic(err)
	}
	return output
}

func justReadGameOutput(ws *websocket.Conn) *conway.Generation {
	var output *conway.Generation
	if err := ws.ReadJSON(&output); err != nil {
		panic(err)
	}
	return output
}

func assertGenerationOutput(ws *websocket.Conn) {
	var output *conway.Generation
	if err := ws.ReadJSON(&output); err != nil {
		Fail("Expected generation output")
	}
}

func assertGenerationTwo(ws *websocket.Conn) {
	var output *conway.Generation
	if err := ws.ReadJSON(&output); err != nil {
		Fail("Expected generation output")
	}

	secondGeneration := &conway.Generation{
		{Point: conway.Point{Row: 2, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 3}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 3, Col: 3}, State: conway.Live, Player: conway.Player1},
	}

	Expect(output).To(Equal(secondGeneration))
}

func sendAckMessage(ws *websocket.Conn, msg string) {
	if err := ws.WriteJSON(map[string]string{"acknowledged": msg}); err != nil {
		panic(err)
	}
}

func httpToWs(u string) string {
	return "ws" + u[len("http"):]
}
