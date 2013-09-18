package game

type State int

const (
  Dead State = iota
  Live
)

type Point struct {
  Row int
  Col int
}

type Cell struct {
  Point
  State State
}

type Generation []Cell

type Game struct {
  Rows int
  Cols int
}

func (this *Game) PointsToGeneration(points *[]Point) *Generation {
  generation := Generation{}
  for _, point := range *points {
    generation = append(generation, Cell{Point: point, State: Live})
  }
  return &generation
}

func (this *Game) GenerationToPoints(generation *Generation) *[]Point {
  points := []Point{}
  for _, cell := range *generation {
    points = append(points, Point{Row: cell.Row, Col: cell.Col})
  }
  return &points
}

func (this *Game) NextGeneration(g *Generation) (next_generation *Generation) {

  next_generation = &Generation{}

  for _, cell := range *g {

    live_cnt := 0
    // Go around cell neighbours
    for _, point := range neighbour_cells_coords(cell.Row, cell.Col) {

      // live neighbour
      if ThereIsCellAtPoint(point, g) {
        live_cnt += 1

      // dead neighbour
      } else {

        // try repopulate dead cell
        // if we haven't done this already
        if !ThereIsCellAtPoint(point, next_generation) {
          live_arount_dead_cnt := 0

          // count live neighbours of dead cell
          for _, arount_dead_point := range neighbour_cells_coords(point[0], point[1]) {

            if ThereIsCellAtPoint(arount_dead_point, g) {
              live_arount_dead_cnt += 1
            }
          }
          if live_arount_dead_cnt == 3 {
            *next_generation = append(*next_generation, Cell{Point: Point{Row: point[0], Col: point[1]}, State: Live})
          }
        }
      }
    }
    if live_cnt == 2 || live_cnt == 3 {
      *next_generation = append(*next_generation, cell)
    }
  }
  return this.DiscardOutOfBoundsCells(next_generation)
}

func (this *Game) DiscardOutOfBoundsCells(next_generation *Generation) *Generation {
  var filtered_generation Generation
  for _, cell := range *next_generation {
    if cell.Row <= this.Rows && cell.Col <= this.Cols {
      filtered_generation = append(filtered_generation, cell)
    }
  }
  return &filtered_generation
}

func ThereIsCellAtPoint(point [2]int, g *Generation) bool {
  res := false
  for _, cell := range *g {
    if cell.Row == point[0] && cell.Col == point[1] {
      res = true
    }
  }
  return res
}

func neighbour_cells_coords(row int, col int) (result *[8][2]int) {
  result = &[8][2]int{
    {row-1, col},
    {row-1, col+1},
    {row, col+1},
    {row+1, col+1},
    {row+1, col},
    {row+1, col-1},
    {row, col-1},
    {row-1, col-1},
  }

  return result
}
