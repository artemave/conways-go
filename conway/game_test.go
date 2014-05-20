package conway

import (
	"sort"
	"testing"
)

var GliderSeed = Generation{
	{Point: Point{Row: 2, Col: 1}, State: Live},
	{Point: Point{Row: 2, Col: 2}, State: Live},
	{Point: Point{Row: 2, Col: 3}, State: Live},
	{Point: Point{Row: 1, Col: 3}, State: Live},
	{Point: Point{Row: 0, Col: 2}, State: Live},
}

var StickSeed = Generation{
	{Point: Point{Row: 0, Col: 0}, State: Live},
	{Point: Point{Row: 0, Col: 1}, State: Live},
	{Point: Point{Row: 0, Col: 2}, State: Live},
}

type ByGeneration Generation

func (g ByGeneration) Len() int {
	return len(g)
}
func (g ByGeneration) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}
func (g ByGeneration) Less(i, j int) bool {
	if g[i].Row == g[j].Row {
		return g[i].Col < g[j].Col
	}
	return g[i].Row < g[j].Row
}

func (this Generation) equal(other Generation) bool {
	if len(this) != len(other) {
		return false
	}

	sorted_this := this
	sorted_other := other

	sort.Sort(ByGeneration(sorted_this))
	sort.Sort(ByGeneration(sorted_other))

	for i, cell := range sorted_this {
		if cell != sorted_other[i] {
			return false
		}
	}
	return true
}

func TestGliderGen1(t *testing.T) {

	ExpectedGen := Generation{
		Cell{Point: Point{Row: 3, Col: 2}, State: Live},
		Cell{Point: Point{Row: 2, Col: 2}, State: Live},
		Cell{Point: Point{Row: 2, Col: 3}, State: Live},
		Cell{Point: Point{Row: 1, Col: 1}, State: Live},
		Cell{Point: Point{Row: 1, Col: 3}, State: Live},
	}

	game := &Game{20, 20}

	ActualGen := game.NextGeneration(&GliderSeed)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestGliderGen2(t *testing.T) {

	ExpectedGen := Generation{
		Cell{Point: Point{Row: 3, Col: 3}, State: Live},
		Cell{Point: Point{Row: 3, Col: 2}, State: Live},
		Cell{Point: Point{Row: 2, Col: 3}, State: Live},
		Cell{Point: Point{Row: 2, Col: 1}, State: Live},
		Cell{Point: Point{Row: 1, Col: 3}, State: Live},
	}

	game := &Game{20, 20}

	ActualGen := game.NextGeneration(&GliderSeed)
	ActualGen = game.NextGeneration(ActualGen)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestStillLife(t *testing.T) {
	ExpectedGen := Generation{
		Cell{Point: Point{Row: 1, Col: 1}, State: Live},
		Cell{Point: Point{Row: 1, Col: 2}, State: Live},
		Cell{Point: Point{Row: 2, Col: 1}, State: Live},
		Cell{Point: Point{Row: 2, Col: 2}, State: Live},
	}

	game := &Game{20, 20}

	ActualGen := game.NextGeneration(&ExpectedGen)
	ActualGen = game.NextGeneration(ActualGen)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestCellsDontSpreadBeyondGameBoundaries(t *testing.T) {
	ExpectedGen := Generation{
		//{Point: Point{Row: 3, Col: 2}, State: Live},
		{Point: Point{Row: 2, Col: 2}, State: Live},
		//{Point: Point{Row: 2, Col: 3}, State: Live},
		{Point: Point{Row: 1, Col: 1}, State: Live},
		//{Point: Point{Row: 1, Col: 3}, State: Live},
	}

	game := &Game{2, 2}

	ActualGen := game.NextGeneration(&GliderSeed)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestCellsDontSpreadBeyondGameBoundariesNegative(t *testing.T) {
	ExpectedGen := Generation{
		/* {Point: Point{Row: -1, Col: 1}, State: Live}, */
		{Point: Point{Row: 0, Col: 1}, State: Live},
		{Point: Point{Row: 1, Col: 1}, State: Live},
	}

	game := &Game{2, 3}

	ActualGen := game.NextGeneration(&StickSeed)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestRepopulatedCellRetainsPlayer(t *testing.T) {
	stick := Generation{
		{Point: Point{Row: 1, Col: 1}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 3}, State: Live, Player: Player1},
	}

	ExpectedGen := Generation{
		{Point: Point{Row: 0, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 2, Col: 2}, State: Live, Player: Player1},
	}

	game := &Game{5, 5}
	ActualGen := game.NextGeneration(&stick)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestCellsBecomeNeutralIfPlayersMixed(t *testing.T) {
	stick := Generation{
		{Point: Point{Row: 1, Col: 1}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 3}, State: Live, Player: Player2},
	}

	ExpectedGen := Generation{
		{Point: Point{Row: 0, Col: 2}, State: Live, Player: None},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: None},
		{Point: Point{Row: 2, Col: 2}, State: Live, Player: None},
	}

	game := &Game{5, 5}
	ActualGen := game.NextGeneration(&stick)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestResurrectedCellsAcquirePlayer(t *testing.T) {
	stick := Generation{
		{Point: Point{Row: 1, Col: 1}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 3}, State: Live, Player: Player1},
	}

	ExpectedGen := Generation{
		{Point: Point{Row: 0, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 2, Col: 2}, State: Live, Player: Player1},
	}

	game := &Game{5, 5}
	ActualGen := game.NextGeneration(&stick)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestResurrectedCellsStayNeutral(t *testing.T) {
	stick := Generation{
		{Point: Point{Row: 1, Col: 1}, State: Live, Player: None},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: None},
		{Point: Point{Row: 1, Col: 3}, State: Live, Player: None},
	}

	ExpectedGen := Generation{
		{Point: Point{Row: 0, Col: 2}, State: Live, Player: None},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: None},
		{Point: Point{Row: 2, Col: 2}, State: Live, Player: None},
	}

	game := &Game{5, 5}
	ActualGen := game.NextGeneration(&stick)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}

func TestNeutralLiveCellRegainedByPlayer(t *testing.T) {
	square := Generation{
		{Point: Point{Row: 1, Col: 1}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: None},
		{Point: Point{Row: 2, Col: 2}, State: Live, Player: None},
		{Point: Point{Row: 2, Col: 1}, State: Live, Player: None},
	}

	ExpectedGen := Generation{
		{Point: Point{Row: 1, Col: 1}, State: Live, Player: Player1},
		{Point: Point{Row: 1, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 2, Col: 2}, State: Live, Player: Player1},
		{Point: Point{Row: 2, Col: 1}, State: Live, Player: Player1},
	}

	game := &Game{5, 5}
	ActualGen := game.NextGeneration(&square)

	if !ActualGen.equal(ExpectedGen) {
		t.Errorf("Actual generation: ", *ActualGen, " is not equal to expected generation: ", ExpectedGen)
	}
}
