require 'es5-shim'
emitEscape = require '../js/emit_escape'
React      = require 'react/addons'
ButtonBar  = React.createFactory(require '../js/button_bar')
TestUtils  = React.addons.TestUtils

describe "ButtonBar"
  describe "when button clicked"
    it "triggers event (to grid) about which button is clicked" @(done)
      bb = TestUtils.renderIntoDocument(ButtonBar {show = true, player = 1})
      buttonLine = TestUtils.findRenderedDOMComponentWithClass(bb, 'button shape line player1')

      cb(e) =
        e.stopPropagation()
        expect(e.detail.shape).to.eq 'line'
        done()

      window.eventServer.once('about-to-place-shape', cb)

      TestUtils.Simulate.click(buttonLine)


  describe "when escape is pressed"
    it "triggers event (to grid) that no button is currently clicked" @(done)
      bb = TestUtils.renderIntoDocument(ButtonBar {show = true})

      cb(e) =
        done()

      window.eventServer.once('no-shape-wants-to-be-placed', cb)

      emitEscape()
