package main

import "github.com/artemave/conways-go/game"

var TestGameRepo = &gamesRepo
var TestDelay = &game.Delay

func (gr *InMemoryGamesRepo) Empty() {
	gr.Games = []*game.Game{}
}
