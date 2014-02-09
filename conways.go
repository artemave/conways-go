package main

import (
	"github.com/araddon/gou"
	"github.com/artemave/conways-go/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

	routes.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	gou.Debug("listening at " + port)
	http.ListenAndServe(":"+port, nil)
}
