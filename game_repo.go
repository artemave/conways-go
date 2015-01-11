package main

import (
	"errors"

	"github.com/artemave/conways-go/conway"
	. "github.com/artemave/conways-go/game"
)

type GamesRepo struct {
	Games []*Game
}

func NewGamesRepo() *GamesRepo {
	gr := &GamesRepo{
		Games: []*Game{},
	}
	return gr
}

func (gr *GamesRepo) FindGameById(id string) *Game {
	for _, game := range gr.Games {
		if game.Id == id {
			return game
		}
	}
	return nil
}

func (gr *GamesRepo) CreateGameById(id string, gameSize string, startGeneration *conway.Generation) (*Game, error) {
	for _, game := range gr.Games {
		if game.Id == id {
			return nil, errors.New("Game with id '" + id + "' is already created.")
		}
	}
	newGame := NewGame(id, gameSize, startGeneration)
	gr.Games = append(gr.Games, newGame)
	return newGame, nil
}

func (gr *GamesRepo) CreatePracticeGameById(id string, startGeneration *conway.Generation) (*Game, error) {
	newGame, err := gr.CreateGameById(id, "small", startGeneration)

	if err != nil {
		return nil, err
	}

	newGame.IsPractice = true
	return newGame, nil
}
