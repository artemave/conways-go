package main

import (
	"log"
	"net/http"
	"os"

	"github.com/araddon/gou"
)

func main() {
	gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

	RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	gou.Debug("listening at " + port)
	http.ListenAndServe(":"+port, nil)
}
