package routes

import (
	"github.com/artemave/conways-go/dependencies/gouuid"
	"net/http"
)

func StartNewGameHandler(w http.ResponseWriter, req *http.Request) {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	http.Redirect(w, req, "/games/"+u.String(), 302)
}

func NewGameHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "../public/index.html")
}
