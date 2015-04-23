package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"code.google.com/p/go-uuid/uuid"

	"fmt"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/config"
	"github.com/artemave/conways-go/conway"
	"github.com/artemave/conways-go/game"
	gga "github.com/artemave/conways-go/google_games_adapter"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var oauthConf = &oauth2.Config{
	ClientID:     config.GoogleClientID(),
	ClientSecret: config.GoogleClientSecret(),
	RedirectURL:  config.OauthRedirectURL(),
	Scopes: []string{
		gga.GamesScope,
	},
	Endpoint: google.Endpoint,
}

var gamesRepo = NewGamesRepo()

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

// RegisterRoutes : registers routes
func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/games", createGameHandler).Methods("POST")
	r.HandleFunc("/practice", createPracticeGameHandler).Methods("POST")
	r.HandleFunc("/submit_score", submitScoreHandler)
	r.HandleFunc("/oauth2callback", oauthCallbackHander)
	r.HandleFunc("/games/{id}", rootHandler)
	r.HandleFunc("/games/play/{id}", GamePlayHandler)
	r.HandleFunc("/fetch_leaderboards", fetchLeaderboardsHandler)
	r.HandleFunc("/scores", scoresIndexHandler)
	r.HandleFunc("/leaderboards", rootHandler)

	return r
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}

func scoresIndexHandler(w http.ResponseWriter, req *http.Request) {
	session, err := sessionCache.Get(req, "sessionCache")
	if err != nil {
		gou.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if scores := session.Values["scores"]; scores != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(scores.(string)))
	}
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	gameSize := r.PostFormValue("gameSize")

	u4 := uuid.New()
	_, err := gamesRepo.CreateGameById(u4, gameSize, startGeneration[gameSize])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprintf(w, u4)
}

func createPracticeGameHandler(w http.ResponseWriter, r *http.Request) {
	u4 := uuid.New()
	newGame, err := gamesRepo.CreateGameById(u4, "small", practiceGameStartGeneration)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	newGame.IsPractice = true
	fmt.Fprintf(w, u4)
}

func submitScoreHandler(w http.ResponseWriter, req *http.Request) {
	// TODO implement real CSRF protection instead of "state"
	gameID := req.URL.Query().Get("gameID")
	state := map[string]string{"gameID": gameID, "callbackFor": "submit_score"}
	stateJSON, _ := json.Marshal(state)
	url := oauthConf.AuthCodeURL(string(stateJSON))
	http.Redirect(w, req, url, 302)
}

func fetchLeaderboardsHandler(w http.ResponseWriter, req *http.Request) {
	state := map[string]string{"callbackFor": "fetch_leaderboards"}
	stateJSON, _ := json.Marshal(state)
	url := oauthConf.AuthCodeURL(string(stateJSON))
	http.Redirect(w, req, url, 302)
}

func oauthCallbackHander(w http.ResponseWriter, req *http.Request) {
	stateJSON := req.URL.Query().Get("state")
	var state map[string]string

	if err := json.Unmarshal(
		[]byte(stateJSON), &state); err != nil {

		gou.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	code := req.URL.Query().Get("code")
	tok, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		gou.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := oauthConf.Client(oauth2.NoContext, tok)

	gapi, err := gga.New(client)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		gou.Error(err)
		return
	}

	session, err := sessionCache.Get(req, "sessionCache")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		gou.Error(err)
		return
	}

	switch state["callbackFor"] {
	case "submit_score":
		if err := processSubmitScore(gapi, state); err != nil {
			gou.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "fetch_leaderboards":
		if err := processFetchLeaderboards(session, gapi); err != nil {
			gou.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		session.Save(req, w)
	default:
		gou.Error(fmt.Printf("callbackFrom is not set: %#v", state))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/leaderboards", 302)
}

func processFetchLeaderboards(session *sessions.Session, gapi interface {
	Leaderboards() ([]gga.Leaderboard, error)
	Scores(gga.Leaderboard) (*gga.LeaderboardScores, error)
}) error {

	leaderboards, err := gapi.Leaderboards()
	if err != nil {
		return err
	}

	var boardsScores []*gga.LeaderboardScores

	var wg sync.WaitGroup
	res := make(chan *gga.LeaderboardScores, len(leaderboards))
	errs := make(chan error, len(leaderboards))

	defer close(errs)
	defer close(res)

	wg.Add(len(leaderboards))

	for _, board := range leaderboards {
		go func(b gga.Leaderboard) {
			defer wg.Done()

			scores, err := gapi.Scores(b)
			if err != nil {
				errs <- err
			}
			res <- scores
		}(board)
	}
	wg.Wait()

FORZ:
	for {
		select {
		case err := <-errs:
			return err
		case scores := <-res:
			boardsScores = append(boardsScores, scores)
		default:
			break FORZ
		}
	}

	session.Options = &sessions.Options{
		MaxAge: 1800,
	}
	cache, _ := json.Marshal(boardsScores)
	session.Values["scores"] = string(cache)

	return nil
}

func processSubmitScore(gapi interface {
	Leaderboards() ([]gga.Leaderboard, error)
	CurrentPlayerScore(gga.Leaderboard) (*gga.PlayerScore, error)
	SubmitScore(gga.Leaderboard, int64) error
}, state map[string]string) error {
	game, err := validateGameSubmitScore(state["gameID"])
	if err != nil {
		return err
	}

	leaderboards, err := gapi.Leaderboards()
	if err != nil {
		return err
	}

	for _, board := range leaderboards {
		if board.Name == game.Size {
			score, err := gapi.CurrentPlayerScore(board)
			if err != nil {
				return err
			}

			err = gapi.SubmitScore(board, score.ScoreValue+3)
			if err != nil {
				return err
			}
			game.SetScoredBy(score.PlayerId)
			break
		}
	}

	return nil
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
