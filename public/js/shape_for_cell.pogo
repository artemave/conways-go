Shape = prototype {
  matrix = math.matrix()

  flipAcrossYeqX() =
    self

  cells(center cell) =
    self.points().map @(point)
      {
        Col = center cell.Col+point.0
        Row = center cell.Row+point.1
      }

  points() =
    math.transpose(self.matrix)._data.map @(coord)
      [coord.0, coord.1 * -1]
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
    self
}

shapeOf(shape) =
  when (shape) [
    is 'line'
      Line()

    is 'square'
      Square()

    is 'glider'
      Glider()
  ]

module.exports = shapeOf
