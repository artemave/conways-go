package routes_test

import (
  . "github.com/artemave/conways-go/dependencies/ginkgo"
  . "github.com/artemave/conways-go/dependencies/gomega"
  "github.com/gorilla/websocket"
  "net/http/httptest"
)

var _ = Describe("GameHandshakeHandler", func() {

  Context("New game", func() {
    It("tells web client to wait", func() {
      ws := wsRequest("/games/handshake/122")
      defer ws.Close()

      output := justRead(ws)
      Expect(output["handshake"]).To(Equal("wait"))
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
      Expect(output["handshake"]).To(Equal("ready"))
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

func justRead(ws *websocket.Conn) map[string]string {
  var output map[string]string
  err := ws.ReadJSON(&output)
  if err != nil {
    panic(err)
  }
  return output
}

func httpToWs(u string) string {
  return "ws" + u[len("http"):]
}
