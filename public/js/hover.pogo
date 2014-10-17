shapeOf = require './shape_for_cell'

Hover(grid) =
  currently pressed button = nil

  document.addEventListener 'about-to-place-shape' @(e)
    currently pressed button := e.detail.shape

  document.addEventListener 'no-shape-wants-to-be-placed'
    grid.clearClass 'hover'
    currently pressed button := nil

  this.maybeDrawShape(cell) =
    if (!currently pressed button)
      return

    shape = shapeOf(currently pressed button)
    if (grid.player == 2)
      shape.flipAcrossYeqX()

    cells = shape.cells(cell)

    if (grid.any of (cells) classed with any of (['fog', 'live', 'new']))
      grid.clearClass 'hover'
    else
      grid.add class 'hover' to (cells)

  this

module.exports = Hover
