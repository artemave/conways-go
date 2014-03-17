package routes_test

import (
  . "github.com/artemave/conways-go/dependencies/ginkgo"
  . "github.com/artemave/conways-go/dependencies/gomega"
  . "github.com/artemave/conways-go/routes"
  "github.com/gorilla/websocket"
  "net/http/httptest"
)

var _ = Describe("GamePlayHandler", func() {

  BeforeEach(func() {
    TestGameRepo.Empty()
  })

  Context("New game", func() {
    It("tells web client to wait", func() {
      ws := wsRequest("/games/play/122")
      defer ws.Close()

      output := justRead(ws)
      Expect(output["handshake"]).To(Equal("wait"))
    })
  })

  Context("Existing game", func() {
    var firstWs *websocket.Conn

    BeforeEach(func() {
      firstWs = wsRequest("/games/play/123")
      justRead(firstWs)
    })
    AfterEach(func() {
      firstWs.Close()
    })

    Context("second client", func() {
      It("tells all web clients to join the game", func() {
        ws := wsRequest("/games/play/123")
        defer ws.Close()

        output := justRead(ws)
        Expect(output["handshake"]).To(Equal("ready"))
        output = justRead(firstWs)
        Expect(output["handshake"]).To(Equal("ready"))
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

        output := justRead(ws)
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
