shapeForCell(shape, cell) =
  when (shape) [
    is 'line'
      [
        { Row = cell.Row, Col = cell.Col-1 }
        { Row = cell.Row, Col = cell.Col }
        { Row = cell.Row, Col = cell.Col+1 }
      ]

    is 'square'
      [
        { Row = cell.Row, Col = cell.Col }
        { Row = cell.Row, Col = cell.Col+1 }
        { Row = cell.Row+1, Col = cell.Col+1 }
        { Row = cell.Row+1, Col = cell.Col }
      ]

    is 'glider'
      [
        { Row = cell.Row-1, Col = cell.Col }
        { Row = cell.Row, Col = cell.Col }
        { Row = cell.Row+1, Col = cell.Col }
        { Row = cell.Row+2, Col = cell.Col-1 }
        { Row = cell.Row+1, Col = cell.Col-2 }
      ]
  ]

module.exports = shapeForCell
