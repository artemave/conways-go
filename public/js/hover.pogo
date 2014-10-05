require './when'
shapeForCell = require './shape_for_cell'

Hover(grid) =
  currently pressed button = nil

  document.addEventListener 'about-to-place-shape' @(e)
    currently pressed button := e.detail.shape

  document.addEventListener 'no-shape-wants-to-be-placed' @(e)
    grid.clearClass 'hover'
    currently pressed button := nil

  this.maybeDrawShape(cell) =
    if (!currently pressed button)
      return

    cells = shapeForCell(currently pressed button, cell)

    if (grid.any of (cells) classed with any of (['fog', 'live']))
      grid.clearClass 'hover'
    else
      grid.add class 'hover' to (cells)

  this

module.exports = Hover
