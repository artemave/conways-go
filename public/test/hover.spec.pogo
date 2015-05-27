Hover       = require '../js/hover'
eventServer = require '../js/event_server'
shapeOf     = require '../js/shape_for_cell'

describe "Hover"

  ['line', 'square', 'glider'].forEach @(s)
    describe "when #(s) button is pressed"
      grid = {}
      hover = nil

      beforeEach
        grid.addClassTo = sinon.spy()
        grid.any of classed with any of () =
          false

        hover := @new Hover(grid)
        eventServer.emit "about-to-place-shape" {detail = {shape = s}}

      it "casts a hover (of #(s)) onto the grid"
        cell = {Row = 2, Col = 2}

        hover.maybeDrawShape(cell)
        expect(grid.add class to).to.have.been.calledWith('hover', shapeOf(s).cells(cell))
