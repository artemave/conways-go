package game

import (
  "testing"
)

var GliderSeed = Generation{
  &Cell{Row: 2, Col: 2},
  &Cell{Row: 2, Col: 3},
  &Cell{Row: 2, Col: 4},
  &Cell{Row: 1, Col: 4},
  &Cell{Row: 0, Col: 3},
}

func TestGliderGen1(t *testing.T) {

  ExpectedGen := Generation{
    &Cell{Row: 3, Col: 3},
    &Cell{Row: 2, Col: 3},
    &Cell{Row: 2, Col: 4},
    &Cell{Row: 1, Col: 2},
    &Cell{Row: 1, Col: 4},
  }

  game := &Game{20, 20}

  ActualGen := game.NextGeneration(&GliderSeed)

  AssertActualGenMatchesExpected(t, ActualGen, &ExpectedGen)
}

func AssertActualGenMatchesExpected(t *testing.T, ActualGen *Generation, ExpectedGen *Generation) {
  t.Errorf("Boom!")
}

