package main

import "github.com/artemave/conways-go/game"

var TestGameRepo = &gamesRepo

func (self *GamesRepo) Empty() {
	self.Games = []*game.Game{}
}

var TestDelay = &game.Delay
