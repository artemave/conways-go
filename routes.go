package main

import (
	"encoding/json"
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"fmt"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/config"
	"github.com/artemave/conways-go/conway"
	gga "github.com/artemave/conways-go/google_games_adapter"
	s "github.com/artemave/conways-go/scores"
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
	_, err := gamesRepo.CreateGameById(u4, gameSize, startGeneration[gameSize], false)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Fprintf(w, u4)
}

func createPracticeGameHandler(w http.ResponseWriter, r *http.Request) {
	u4 := uuid.New()
	_, err := gamesRepo.CreateGameById(u4, "small", practiceGameStartGeneration, true)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
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
		game := gamesRepo.FindGameById(state["gameID"])
		if game == nil {
			gou.Error(fmt.Printf("Game with id '%s' not found", state["gameID"]))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := s.SubmitScore(gapi, game); err != nil {
			gou.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case "fetch_leaderboards":
		boardsScores, err := s.FetchLeaderboards(gapi)
		if err != nil {
			gou.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session.Options = &sessions.Options{
			MaxAge: 1800,
		}
		cache, _ := json.Marshal(boardsScores)
		session.Values["scores"] = string(cache)

		session.Save(req, w)
	default:
		gou.Error(fmt.Printf("callbackFrom is not set: %#v", state))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, req, "/leaderboards", 302)
}
