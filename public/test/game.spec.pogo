require 'es5-shim'
React             = require 'react/addons'
StubRouterContext = require './stub_router_context'
Game              = React.createFactory(require '../js/views/game')
TestUtils         = React.addons.TestUtils

describe 'Game'
  describe 'initialisation'
    wsSpy = null

    beforeEach
      wsSpy := sinon.spy()
      sinon.stub(window, 'WrapWebSocket', @{ wsSpy })

    it 'opens a websocket connection'
      gameComponent = StubRouterContext(
        Game
        { wsHost = 'ws://host' }
        { getCurrentParams() = { gameId = '123' } }
      )

      g = TestUtils.renderIntoDocument(React.createElement(gameComponent))

      expect(wsSpy).to.have.been.calledWithNew
      expect(wsSpy).to.have.been.calledWith('ws://host/games/play/123')

  describe 'responding to ws messages'
    context 'server sent "wait"'
      it 'acknowlekdges the message'
