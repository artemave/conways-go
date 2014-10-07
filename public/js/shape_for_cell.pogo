Shape = prototype {
  matrix = math.matrix()

  flipAcrossYeqX() =
    self.matrix

  cells(center cell) =
    math.transpose(self.matrix)._data.map @(coord)
      {
        Col = center cell.Col+coord.0
        Row = center cell.Row+(coord.1 * -1)
      }
}

Line() = Shape {
  matrix = math.matrix [[-1,0,1],[0,0,0]]
}

Square() = Shape {
  matrix = math.matrix [[0,1,1,0],[0,0,-1,-1]]
}

Glider() = Shape {
  matrix = math.matrix [[1,1,1,0,-1],[1,0,-1,-1,0]]

  flipAcrossYeqX() =
    self.matrix = math.multiply([[0,1],[1,0]], self.matrix)
}

shapeForCell(shape) =
  when (shape) [
    is 'line'
      Line()

    is 'square'
      Square()

    is 'glider'
      Glider()
  ]

module.exports = shapeForCell
