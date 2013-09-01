package main

import (
  "github.com/conways-go/game"
  "github.com/conways-go/screen"
)

var GliderSeed = game.Generation{
  &game.Cell{Row: 2, Col: 2},
  &game.Cell{Row: 2, Col: 3},
  &game.Cell{Row: 2, Col: 4},
  &game.Cell{Row: 1, Col: 4},
  &game.Cell{Row: 0, Col: 3},
}


func main() {
  g := game.Game{Rows: 10, Cols: 20}
  s := screen.Screen{}
  generation := &GliderSeed

  for {
    s.RenderGeneration(generation)
    generation = g.NextGeneration(generation)
  }
}
