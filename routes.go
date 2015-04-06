package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"fmt"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/config"
	"github.com/artemave/conways-go/conway"
	"github.com/artemave/conways-go/game"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleGames "google.golang.org/api/games/v1"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/games", CreateGameHandler).Methods("POST")
	r.HandleFunc("/practice", CreatePracticeGameHandler).Methods("POST")
	r.HandleFunc("/submit_score", submitScoreHandler)
	r.HandleFunc("/oauth2callback", oauthCallbackHander)
	r.HandleFunc("/games/{id}", ShowGameHandler)
	r.HandleFunc("/games/play/{id}", GamePlayHandler)

	return r
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}

var startGeneration = map[string]*conway.Generation{
	"large": &conway.Generation{
		{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 5, Col: 4}, State: conway.Live, Player: conway.Player1},

		{Point: conway.Point{Row: 64, Col: 95}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 65, Col: 95}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 66, Col: 95}, State: conway.Live, Player: conway.Player2},
	},
	"medium": &conway.Generation{
		{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 5, Col: 4}, State: conway.Live, Player: conway.Player1},

		{Point: conway.Point{Row: 44, Col: 75}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 45, Col: 75}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 46, Col: 75}, State: conway.Live, Player: conway.Player2},
	},
	"small": &conway.Generation{
		{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
		{Point: conway.Point{Row: 5, Col: 4}, State: conway.Live, Player: conway.Player1},

		{Point: conway.Point{Row: 20, Col: 35}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 21, Col: 35}, State: conway.Live, Player: conway.Player2},
		{Point: conway.Point{Row: 22, Col: 35}, State: conway.Live, Player: conway.Player2},
	},
}

var practiceGameStartGeneration = &conway.Generation{
	{Point: conway.Point{Row: 3, Col: 4}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 4, Col: 4}, State: conway.Live, Player: conway.Player1},
	{Point: conway.Point{Row: 5, Col: 4}, State: conway.Live, Player: conway.Player1},

	{Point: conway.Point{Row: 20, Col: 35}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 21, Col: 35}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 22, Col: 35}, State: conway.Live, Player: conway.Player2},

	{Point: conway.Point{Row: 3, Col: 35}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 3, Col: 36}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 4, Col: 35}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 4, Col: 36}, State: conway.Live, Player: conway.Player2},

	{Point: conway.Point{Row: 21, Col: 24}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 21, Col: 25}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 21, Col: 26}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 22, Col: 24}, State: conway.Live, Player: conway.Player2},
	{Point: conway.Point{Row: 23, Col: 25}, State: conway.Live, Player: conway.Player2},
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	gameSize := r.PostFormValue("gameSize")

	u4 := uuid.New()
	_, err := gamesRepo.CreateGameById(u4, gameSize, startGeneration[gameSize])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprintf(w, u4)
}

func CreatePracticeGameHandler(w http.ResponseWriter, r *http.Request) {
	u4 := uuid.New()
	newGame, err := gamesRepo.CreateGameById(u4, "small", practiceGameStartGeneration)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	newGame.IsPractice = true
	fmt.Fprintf(w, u4)
}

func ShowGameHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}

func submitScoreHandler(w http.ResponseWriter, req *http.Request) {
	// TODO implement real CSRF protection instead of "state"
	gameID := req.URL.Query().Get("gameID")
	state := map[string]string{"gameID": gameID}
	stateJSON, _ := json.Marshal(state)
	url := conf.AuthCodeURL(string(stateJSON))
	http.Redirect(w, req, url, 302)
}

var gamesRepo = NewGamesRepo()

func oauthCallbackHander(w http.ResponseWriter, req *http.Request) {
	stateJSON := req.URL.Query().Get("state")
	var state map[string]string

	if err := json.Unmarshal(
		[]byte(stateJSON), &state); err != nil {

		gou.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	game, err := validateGameSubmitScore(state["gameID"])
	if err != nil {
		gou.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	code := req.URL.Query().Get("code")
	tok, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		gou.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	gapi, err := googleGames.New(client)
	if err != nil {
		gou.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	leaderboards, err := gapi.Leaderboards.List().Do()
	if err != nil {
		gou.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	//TODO test
	for _, board := range leaderboards.Items {
		if board.Name == game.Size {
			score, err := gapi.Scores.Get("me", board.Id, "ALL_TIME").Do()
			if err != nil {
				gou.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			var newScore int64
			if len(score.Items) == 0 {
				newScore = 3
			} else {
				currentScore := score.Items[0].ScoreValue
				newScore = int64(currentScore) + 3
			}

			res, err := gapi.Scores.Submit(board.Id, newScore).Do()
			if err != nil {
				gou.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			game.SetScoredBy(score.Player.PlayerId)

			gou.Info(fmt.Sprintf("Score submitted: %#v", res))
			break
		}
	}

	http.Redirect(w, req, "/", 302)
}

var conf = &oauth2.Config{
	ClientID:     config.GoogleClientID(),
	ClientSecret: config.GoogleClientSecret(),
	RedirectURL:  config.OauthRedirectURL(),
	Scopes: []string{
		googleGames.GamesScope,
	},
	Endpoint: google.Endpoint,
}

//TODO test
func validateGameSubmitScore(gameID string) (*game.Game, error) {
	game := gamesRepo.FindGameById(gameID)

	if game == nil {
		return nil, errors.New("Game not found: " + gameID)
	}

	if game.IsPractice {
		return nil, errors.New("Score can not be submitted for practice game.")
	}

	if !game.IsFinished {
		return nil, errors.New("Can't submit score for a game that has not yet finished.")
	}

	if game.GetScoredBy() != nil {
		return nil, errors.New("Score has already been submitted.")
	}

	return game, nil
}
