require 'es5-shim'
emitEscape  = require '../js/emit_escape'
eventServer = require '../js/event_server'
React       = require 'react/addons'
ButtonBar   = React.createFactory(require '../js/button_bar')
TestUtils   = React.addons.TestUtils

describe "ButtonBar"
  describe "when button clicked"
    it "triggers event (to grid) about which button is clicked" @(done)
      bb = TestUtils.renderIntoDocument(ButtonBar {show = true, player = 1})
      buttonLine = TestUtils.findRenderedDOMComponentWithClass(bb, 'button shape line player1')

      cb(e) =
        expect(e.detail.shape).to.eq 'line'
        done()

      eventServer.once('about-to-place-shape', cb)

      TestUtils.Simulate.click(buttonLine)


  describe "when escape is pressed"
    original listeners = []

    beforeEach
      original listeners := eventServer.listeners('no-shape-wants-to-be-placed').slice 0
      eventServer.removeAllListeners('no-shape-wants-to-be-placed')

    afterEach
      [l <- original listeners, eventServer.addListener('no-shape-wants-to-be-placed', l)]

    it "triggers event (to grid) that no button is currently clicked" @(done)
      bb = TestUtils.renderIntoDocument(ButtonBar {show = true})

      cb() =
        done()

      eventServer.on('no-shape-wants-to-be-placed', cb)

      emitEscape()
