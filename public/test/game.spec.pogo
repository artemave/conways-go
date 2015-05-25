require 'es5-shim'
React             = require 'react/addons'
Game              = require '../js/views/game'
StubRouterContext = require './stub_router_context'
TestUtils         = React.addons.TestUtils

describe 'Game'
  wsSpy   = null
  sandbox = null

  gameComponent = StubRouterContext(
    React.createFactory(Game)
    { wsHost = 'ws://host' }
    { getCurrentParams() = { gameId = '123' } }
  )

  beforeEach
    sandbox := sinon.sandbox.create()
    sandbox.stub(window, 'WrapWebSocket', @{ wsSpy })

  afterEach
    sandbox.restore()

  describe 'initialisation'
    beforeEach
      wsSpy := sandbox.spy()
      TestUtils.renderIntoDocument(React.createElement(gameComponent))

    it 'opens a websocket connection'
      expect(wsSpy).to.have.been.calledWithNew
      expect(wsSpy).to.have.been.calledWith 'ws://host/games/play/123'

  describe 'responding to ws messages'
    ws                = null
    componentInstance = null

    beforeEach
      ws := { send = sandbox.spy() }
      wsSpy := @() @{ ws }
      componentInstance := TestUtils.renderIntoDocument(React.createElement(gameComponent))

    context 'server sent "wait"'
      beforeEach
        ws.onmessage({ data = JSON.stringify { Handshake = 'wait' }})

      it 'acknowlekdges the message'
        expect(ws.send).to.have.been.calledWith(JSON.stringify { "acknowledged" = "wait" })

      it 'shows share instructions'

      it 'hides game grid (in case "wait" was the result of other player disconnect whilst playing)'
      it 'hides "game is paused" (in case "wait" was the result of other player disconnect whilst game was paused by them)'
