require 'es5-shim'
React             = require 'react/addons'
Game              = require '../js/views/game'
StubRouterContext = require './stub_router_context'
ShareInstructions = require '../js/share_instructions'
GameIsPaused      = require '../js/game_is_paused'
ButtonBar         = require '../js/button_bar'
Grid              = require '../js/grid'
HelpPopup         = require '../js/help_popup'
SubmitScorePopup  = require '../js/submit_score_popup'
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
    gameComponentInstance = null

    beforeEach
      ws := { send = sandbox.spy() }
      wsSpy := @() @{ ws }
      gameComponentInstance := TestUtils.renderIntoDocument(React.createElement(gameComponent))

    context 'server sent "wait"'
      beforeEach
        ws.onmessage({ data = JSON.stringify { Handshake = 'wait' }})

      it 'acknowlekdges the message'
        expect(ws.send).to.have.been.calledWith(JSON.stringify { "acknowledged" = "wait" })

      it 'shows share instructions'
        TestUtils.findRenderedComponentWithType(gameComponentInstance, ShareInstructions)

      ['ButtonBar', 'GameIsPaused', 'Grid', 'HelpPopup', 'SubmitScorePopup'].forEach @(type)
        it "does NOT show #(type)"
          expect(TestUtils.scryRenderedComponentsWithType(gameComponentInstance, eval(type)).length).to.eq 0
