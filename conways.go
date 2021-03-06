package main

import (
	"log"
	"os"

	"github.com/araddon/gou"
	"github.com/artemave/conways-go/config"
	"github.com/codegangsta/negroni"
	"gopkg.in/boj/redistore.v1"
)

var sessionCache *redistore.RediStore

func main() {
	gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}
	n := negroni.Classic()

	// TODO real secret maybe
	r, err := redistore.NewRediStore(10, "tcp", config.RedisHost(), config.RedisPassword(), []byte("secret123"))
	if err != nil {
		panic(err)
	}
	sessionCache = r
	defer sessionCache.Close()

	mux := RegisterRoutes()
	n.UseHandler(mux)
	n.Run(":" + port)
}
