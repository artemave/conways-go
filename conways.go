package main

import (
	"log"
	"os"

	"github.com/araddon/gou"
	"github.com/codegangsta/negroni"
)

func main() {
	gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	n := negroni.Classic()
	mux := RegisterRoutes()

	n.UseHandler(mux)
	n.Run(":" + port)
}
