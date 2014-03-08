package routes_test

import (
  "code.google.com/p/go.net/websocket"
  . "github.com/artemave/conways-go/dependencies/ginkgo"
  . "github.com/artemave/conways-go/dependencies/gomega"
  "net/http/httptest"
  "net/url"
)

var _ = Describe("GameHandshakeHandler", func() {

  Context("New game", func() {
    It("tells web client to wait", func() {
      ws := wsRequest("/games/handshake/123")
      defer ws.Close()

      output := justRead(ws)
      Expect(output).To(Equal("{\"handshake\":\"wait\"}"))
    })
  })

  Context("Existing game", func() {
    var firstWs *websocket.Conn

    BeforeEach(func() {
      firstWs = wsRequest("/games/handshake/123")
    })
    AfterEach(func() {
      firstWs.Close()
    })

    It("tells all web clients to join the game", func() {
      ws := wsRequest("/games/handshake/123")
      defer ws.Close()

      output := justRead(ws)
      Expect(output).To(Equal("{\"handshake\":\"ready\"}"))
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

func justRead(ws *websocket.Conn) string {
  msg := make([]byte, 512)
  n, err := ws.Read(msg)
  if err != nil {
    panic("Read: " + err.Error())
  }
  return string(msg[0:n])
}
