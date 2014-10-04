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

    grid.add class 'hover' to (shapeForCell(currently pressed button, cell))

  this

module.exports = Hover
