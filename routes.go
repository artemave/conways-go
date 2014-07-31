package main

import (
	"net/http"

	"code.google.com/p/go-uuid/uuid"

	"github.com/araddon/gou"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("no-secrets-here--just-storing-state-between-requests"))

func RegisterRoutes() {
	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/games", CreateGameHandler).Methods("POST")
	r.HandleFunc("/games/{id}", ShowGameHandler)
	r.HandleFunc("/games/play/{id}", GamePlayHandler)

	http.Handle("/", r)
	http.Handle("/public/", http.FileServer(http.Dir("./")))
}

func RootHandler(w http.ResponseWriter, req *http.Request) {
	gou.Debug("GET: /")
	http.ServeFile(w, req, "./public/index.html")
}

func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	gou.Debug("POST: /games/")
	gameSize := r.PostFormValue("size")

	session, _ := store.Get(r, "session")
	session.Values["gameSize"] = gameSize
	session.Save(r, w)

	u4 := uuid.New()
	http.Redirect(w, r, "/games/"+u4, 302)
}

func ShowGameHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	gou.Debug("GET: /games/" + id)
	http.ServeFile(w, req, "./public/index.html")
}
