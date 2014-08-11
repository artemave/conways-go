package main

var TestGameRepo = &gamesRepo
var TestStartGeneration = &startGeneration

func (self *GamesRepo) Empty() {
	self.Games = []*Game{}
}

var TestDelay = &delay
