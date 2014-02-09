package routes_test

import (
	. "github.com/artemave/conways-go/dependencies/ginkgo"
	. "github.com/artemave/conways-go/dependencies/gomega"
	. "github.com/artemave/conways-go/routes"
	// "io/ioutil"
	/* "net/http" */
	"code.google.com/p/go.net/websocket"
	"net/http/httptest"
	"net/url"
)

var _ = Describe("GameHandshakeHandler", func() {
	RegisterRoutes()

	Context("New game", func() {
		It("tells web client to wait", func() {
			ws := wsRequest("/handshake/games/123")
			defer ws.Close()

			output := justRead(ws)
			Expect(output).To(Equal("{\"handshake\": \"wait\"}"))
		})
	})

	Context("Existing game", func() {
		It("tells all web clients to join the game", func() {
		})
	})
})

func wsRequest(path string) *websocket.Conn {
	server := httptest.NewServer(nil)
	u, err := url.Parse(server.URL)

	if err != nil {
		panic(err)
	}

	ws, err := websocket.Dial("ws://"+u.Host+path, "", server.URL)
	if err != nil {
		panic(err)
	}
	return ws
}

func justRead(ws *websocket.Conn) []byte {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		panic("Read: " + err.Error())
	}
	return msg[0:n]
}
