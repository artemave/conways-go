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

func (this *Game) NextGeneration(g *Generation) (next_generation *Generation) {

	next_generation = &Generation{}

	for _, cell := range *g {

		live_neighbours := []*Cell{}

		// Go around cell neighbours
		for _, point := range neighbour_cells_coords(cell.Row, cell.Col) {

			// live neighbour
			if lc := CellAtPoint(point, g); lc != nil {
				live_neighbours = append(live_neighbours, lc)

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
						playersAround := []Player{}

						for _, lc := range live_arount_dead {
							if lc.Player != None {
								playersAround = append(playersAround, lc.Player)
							}
						}
						playersAround = uniq_players(playersAround)

						player := None
						if len(playersAround) == 1 {
							player = playersAround[0]
						}

						*next_generation = append(*next_generation, Cell{Point: Point{Row: point[0], Col: point[1]}, State: Live, Player: player})
					}
				}
			}
		}

		if len(live_neighbours) == 2 || len(live_neighbours) == 3 {
			new_cell := cell
			players := []Player{}

			for _, ln := range live_neighbours {
				if ln.Player != None {
					players = append(players, ln.Player)
				}
			}
			players = uniq_players(players)

			if cell.Player == None {
				// only one player around - he gets to own this cell
				if len(players) == 1 {
					new_cell.Player = players[0]
				}
			} else {
				for _, p := range players {
					// any other players around - cell becomes neutral
					if p != cell.Player {
						new_cell.Player = None
					}
				}
			}
			*next_generation = append(*next_generation, new_cell)
		}
	}
	return this.DiscardOutOfBoundsCells(next_generation)
}

func uniq_players(players []Player) []Player {
	m := make(map[Player]bool)
	for _, p := range players {
		m[p] = true
	}
	res := []Player{}

	for k, _ := range m {
		res = append(res, k)
	}
	return res
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
