require 'es5-shim'
React             = require 'react/addons'
StubRouterContext = require './stub_router_context'
Game              = React.createFactory(require '../js/views/game')
TestUtils         = React.addons.TestUtils

describe 'Game'
  wsSpy   = null
  sandbox = null
  gameComponent = StubRouterContext(
    Game
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
    ws = null

    beforeEach
      ws := { send = sandbox.spy() }
      wsSpy := @() @{ ws }
      TestUtils.renderIntoDocument(React.createElement(gameComponent))

    context 'server sent "wait"'
      it 'acknowlekdges the message'
        ws.onmessage({ data = JSON.stringify { Handshake = 'wait' }})
        expect(ws.send).to.have.been.calledWith(JSON.stringify { "acknowledged" = "wait" })
