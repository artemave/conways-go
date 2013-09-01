package game

type Cell struct {
  Row int
  Col int
}

type Generation []*Cell

type Game struct {
  Rows int
  Cols int
}

func (this *Game) NextGeneration(g *Generation) *Generation {
  return g
}
