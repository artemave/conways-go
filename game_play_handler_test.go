package main_test

import (
	"net/http/httptest"
	"time"
	. "github.com/artemave/conways-go"
	"github.com/artemave/conways-go/conway"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GamePlayHandler", func() {

	BeforeEach(func() {
		TestGameRepo.Empty()
	})

	Context("New game", func() {
		It("tells web client to wait", func() {
			ws := wsRequest("/games/play/122")
			defer ws.Close()

			output := justReadHandshake(ws)
			Expect(output["handshake"]).To(Equal("wait"))
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

		Context("second client", func() {
			BeforeEach(func() {
				secondWs = wsRequest("/games/play/123")
			})
			AfterEach(func() {
				secondWs.Close()
			})

			It("tells all web clients to join the game", func() {
				output := justReadHandshake(secondWs)
				Expect(output["handshake"]).To(Equal("ready"))
				output = justReadHandshake(firstWs)
				Expect(output["handshake"]).To(Equal("ready"))
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
						time.Sleep(time.Millisecond * 20)
						secondWs.Close()
					})

					It("tells first client to wait", func() {
						justReadGameOutput(firstWs)

						output := justReadHandshake(firstWs)
						Expect(output["handshake"]).To(Equal("wait"))
					})
				})
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
				Expect(output["handshake"]).To(Equal("game_taken"))
			})
		})

	})
})

var server = httptest.NewServer(nil)

func wsRequest(path string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(httpToWs(server.URL+path), nil)
	if err != nil {
		panic("Dial() returned error " + err.Error())
	}
	return ws
}

func justReadHandshake(ws *websocket.Conn) map[string]string {
	var output map[string]string
	err := ws.ReadJSON(&output)
	if err != nil {
		panic(err)
	}
	return output
}

func justReadGameOutput(ws *websocket.Conn) *[]conway.Point {
	var output *[]conway.Point
	if err := ws.ReadJSON(&output); err != nil {
		panic(err)
	}
	return output
}

func assertGenerationOutput(ws *websocket.Conn) {
	var output *[]conway.Point
	if err := ws.ReadJSON(&output); err != nil {
		Fail("Expected generation output")
	}
}

func sendAckMessage(ws *websocket.Conn, msg string) {
	if err := ws.WriteJSON(map[string]string{"acknowledged": msg}); err != nil {
		panic(err)
	}
}

func httpToWs(u string) string {
	return "ws" + u[len("http"):]
}
