ButtonBar = require '../js/button_bar'
TestUtils = React.addons.TestUtils

describe "ButtonBar"
  describe "when button clicked"
    it "triggers event (to grid) about which button is clicked" @(done)
      bb = @new ButtonBar()
      comp = TestUtils.renderIntoDocument(bb.reactComponent())
      buttonLine = TestUtils.findRenderedDOMComponentWithClass(comp, 'button line')

      cb(e) =
        expect(e.detail.shape).to.eq 'line'
        done()

      document.addEventListener('about-to-place-shape', cb)

      TestUtils.Simulate.click(buttonLine)
      

  describe "when escape is pressed"
    it "triggers event (to grid) that no button is currently clicked" @(done)
      bb = @new ButtonBar()
      comp = TestUtils.renderIntoDocument(bb.reactComponent())

      cb(e) =
        done()

      document.addEventListener('no-shape-wants-to-be-placed', cb)

      e = document.createEvent "Events"
      e.initEvent("keydown", true, true)
      e.keyCode = 27
      e.which = 27
      document.dispatchEvent(e)
