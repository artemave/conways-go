package conway

type State int

const (
	Dead State = iota
	Live
)

type Player int

const (
	None Player = iota
	Player1
	Player2
)

type Point struct {
	Row int
	Col int
}

type Cell struct {
	Point
	State  State
	Player Player
}

type Generation []Cell

type Game struct {
	Rows int
	Cols int
}

func (this *Generation) AddPoints(points []Point, player Player) *Generation {
	for _, point := range points {
		*this = append(*this, Cell{Point: point, State: Live, Player: player})
	}
	return this
}

// func PointsToGeneration(points *[]Point) *Generation {
// 	generation := Generation{}
// 	for _, point := range *points {
// 		generation = append(generation, Cell{Point: point, State: Live})
// 	}
// 	return &generation
// }

func GenerationToPoints(generation *Generation) *[]Point {
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
			if CellAtPoint(point, g) != nil {
				live_cnt += 1

				// dead neighbour
			} else {

				// try repopulate dead cell
				// if we haven't done this already
				if CellAtPoint(point, next_generation) == nil {
					live_arount_dead := []*Cell{}

					// count live neighbours of dead cell
					for _, arount_dead_point := range neighbour_cells_coords(point[0], point[1]) {

						if c := CellAtPoint(arount_dead_point, g); c != nil {
							live_arount_dead = append(live_arount_dead, c)
						}
					}
					if len(live_arount_dead) == 3 {
						player := None
						if live_arount_dead[0].Player == Player1 && live_arount_dead[1].Player == Player1 && live_arount_dead[2].Player == Player1 {
							player = Player1
						} else if live_arount_dead[0].Player == Player2 && live_arount_dead[1].Player == Player2 && live_arount_dead[2].Player == Player2 {
							player = Player2
						} else {
							player = None
						}

						*next_generation = append(*next_generation, Cell{Point: Point{Row: point[0], Col: point[1]}, State: Live, Player: player})
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
		if cell.Row >= 0 && cell.Row <= this.Rows && cell.Col >= 0 && cell.Col <= this.Cols {
			filtered_generation = append(filtered_generation, cell)
		}
	}
	return &filtered_generation
}

func CellAtPoint(point [2]int, g *Generation) *Cell {
	for _, cell := range *g {
		if cell.Row == point[0] && cell.Col == point[1] {
			return &cell
		}
	}
	return nil
}

func neighbour_cells_coords(row int, col int) (result *[8][2]int) {
	result = &[8][2]int{
		{row - 1, col},
		{row - 1, col + 1},
		{row, col + 1},
		{row + 1, col + 1},
		{row + 1, col},
		{row + 1, col - 1},
		{row, col - 1},
		{row - 1, col - 1},
	}

	return result
}

func GosperGliderGun() *Generation {
	GosperGliderGun := Generation{
		{Point: Point{Row: 5, Col: 1}, State: Live},
		{Point: Point{Row: 5, Col: 2}, State: Live},
		{Point: Point{Row: 6, Col: 1}, State: Live},
		{Point: Point{Row: 6, Col: 2}, State: Live},
		{Point: Point{Row: 3, Col: 13}, State: Live},
		{Point: Point{Row: 3, Col: 14}, State: Live},
		{Point: Point{Row: 4, Col: 12}, State: Live},
		{Point: Point{Row: 4, Col: 16}, State: Live},
		{Point: Point{Row: 5, Col: 11}, State: Live},
		{Point: Point{Row: 5, Col: 17}, State: Live},
		{Point: Point{Row: 6, Col: 11}, State: Live},
		{Point: Point{Row: 6, Col: 15}, State: Live},
		{Point: Point{Row: 6, Col: 17}, State: Live},
		{Point: Point{Row: 6, Col: 18}, State: Live},
		{Point: Point{Row: 7, Col: 11}, State: Live},
		{Point: Point{Row: 7, Col: 17}, State: Live},
		{Point: Point{Row: 8, Col: 12}, State: Live},
		{Point: Point{Row: 8, Col: 16}, State: Live},
		{Point: Point{Row: 9, Col: 13}, State: Live},
		{Point: Point{Row: 9, Col: 14}, State: Live},
		{Point: Point{Row: 1, Col: 25}, State: Live},
		{Point: Point{Row: 2, Col: 23}, State: Live},
		{Point: Point{Row: 2, Col: 25}, State: Live},
		{Point: Point{Row: 3, Col: 21}, State: Live},
		{Point: Point{Row: 3, Col: 22}, State: Live},
		{Point: Point{Row: 4, Col: 21}, State: Live},
		{Point: Point{Row: 4, Col: 22}, State: Live},
		{Point: Point{Row: 5, Col: 21}, State: Live},
		{Point: Point{Row: 5, Col: 22}, State: Live},
		{Point: Point{Row: 6, Col: 23}, State: Live},
		{Point: Point{Row: 6, Col: 25}, State: Live},
		{Point: Point{Row: 7, Col: 25}, State: Live},
		{Point: Point{Row: 3, Col: 35}, State: Live},
		{Point: Point{Row: 3, Col: 36}, State: Live},
		{Point: Point{Row: 4, Col: 35}, State: Live},
		{Point: Point{Row: 4, Col: 36}, State: Live},
	}
	return &GosperGliderGun
}
