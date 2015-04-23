package main

import (
	"errors"

	"github.com/artemave/conways-go/conway"
	. "github.com/artemave/conways-go/game"
)

type GamesRepo interface {
	FindGameById(id string) *Game
	CreateGameById(id string, gameSize string, startGeneration *conway.Generation, isPractice bool) (*Game, error)
}

type InMemoryGamesRepo struct {
	Games []*Game
}

func NewGamesRepo() GamesRepo {
	gr := &InMemoryGamesRepo{
		Games: []*Game{},
	}
	return gr
}

func (gr *InMemoryGamesRepo) FindGameById(id string) *Game {
	for _, game := range gr.Games {
		if game.Id == id {
			return game
		}
	}
	return nil
}

func (gr *InMemoryGamesRepo) CreateGameById(id string, gameSize string, startGeneration *conway.Generation, isPractice bool) (*Game, error) {
	for _, game := range gr.Games {
		if game.Id == id {
			return nil, errors.New("Game with id '" + id + "' is already created.")
		}
	}

	newGame := NewGame(id, gameSize, startGeneration, isPractice)
	gr.Games = append(gr.Games, newGame)
	return newGame, nil
}
