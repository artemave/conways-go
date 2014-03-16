package routes

var TestGameRepo = &gamesRepo

func (self *GamesRepo) Empty() {
  self.Games = []*Game{}
}
