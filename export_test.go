package main

var TestGameRepo = &gamesRepo

func (self *GamesRepo) Empty() {
	self.Games = []*Game{}
}

var TestStartGeneration = &startGeneration
var TestDelay = &delay
