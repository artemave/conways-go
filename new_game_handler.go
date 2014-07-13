package main

import (
	"net/http"

	"code.google.com/p/go-uuid/uuid"
	"github.com/araddon/gou"
)

func StartNewGameHandler(w http.ResponseWriter, req *http.Request) {
	gou.Debug("/")

	u4 := uuid.New()
	http.Redirect(w, req, "/games/"+u4, 302)
}

func NewGameHandler(w http.ResponseWriter, req *http.Request) {
	gou.Debug("/games/")

	http.ServeFile(w, req, "./public/index.html")
}
