Hover = require '../js/hover'
shapeOf = require '../js/shape_for_cell'

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

        e = @new CustomEvent "about-to-place-shape" {detail = {shape = s}}
        document.dispatchEvent(e)

      it "casts a hover (of #(s)) onto the grid"
        cell = {Row = 2, Col = 2}

        hover.maybeDrawShape(cell)
        expect(grid.add class to).to.have.been.calledWith('hover', shapeOf(s).cells(cell))
