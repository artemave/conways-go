require 'es5-shim'
emitEscape = require '../js/emit_escape'
ButtonBar  = require '../js/button_bar'
React      = require 'react/addons'
TestUtils  = React.addons.TestUtils

describe "ButtonBar"
  cb = nil

  describe "when button clicked"
    afterEach
      document.removeEventListener('about-to-place-shape', cb)

    it "triggers event (to grid) about which button is clicked" @(done)
      bb = TestUtils.renderIntoDocument(ButtonBar {show = true})
      buttonLine = TestUtils.findRenderedDOMComponentWithClass(bb, 'button line')

      cb(e) :=
        expect(e.detail.shape).to.eq 'line'
        done()

      document.addEventListener('about-to-place-shape', cb)

      TestUtils.Simulate.click(buttonLine)


  describe "when escape is pressed"
    afterEach
      document.removeEventListener('no-shape-wants-to-be-placed', cb)

    it "triggers event (to grid) that no button is currently clicked" @(done)
      bb = TestUtils.renderIntoDocument(ButtonBar {show = true})

      cb(e) :=
        done()

      document.addEventListener('no-shape-wants-to-be-placed', cb)

      emitEscape()
