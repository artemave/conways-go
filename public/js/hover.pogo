require './when'

Hover(svg) =
  currently pressed button = nil

  document.addEventListener 'about-to-place-shape' @(e)
    currently pressed button := e.detail.shape

  document.addEventListener 'no-shape-wants-to-be-placed' @(e)
    svg.selectAll 'rect.hover'.classed 'hover' (false)
    currently pressed button := nil

  this.maybeDrawShape(cell) =
    if (!currently pressed button)
      return

    shape = []
    when(currently pressed button) [
      is 'line'
        shape := [
          { Row = cell.Row, Col = cell.Col-1 }
          { Row = cell.Row, Col = cell.Col }
          { Row = cell.Row, Col = cell.Col+1 }
        ]
    ]

    svg.selectAll 'rect'.data(shape) @(d) @{ "#(d.Row)_#(d.Col)" }.
    classed 'hover' (true).
    exit().classed 'hover' (false)

  this

module.exports = Hover
