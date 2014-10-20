package main

import (
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"fmt"
	"github.com/artemave/conways-go/conway"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"os"
)

func RegisterRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/games", CreateGameHandler).Methods("POST")
	r.HandleFunc("/games/{id}", ShowGameHandler)
	r.HandleFunc("/games/play/{id}", GamePlayHandler)

	http.Handle("/", handlers.LoggingHandler(os.Stdout, r))
	http.Handle("/public/", http.FileServer(http.Dir("./")))
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

func ShowGameHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}
