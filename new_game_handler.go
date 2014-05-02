package main

import (
	"net/http"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/dependencies/gouuid"
)

func StartNewGameHandler(w http.ResponseWriter, req *http.Request) {
	gou.Debug("/")

	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	http.Redirect(w, req, "/games/"+u.String(), 302)
}

func NewGameHandler(w http.ResponseWriter, req *http.Request) {
	gou.Debug("/games/")

	http.ServeFile(w, req, "./public/index.html")
}
