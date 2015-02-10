package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/araddon/gou"
)

func main() {
	defer RecoverAndLogError()

	gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

	RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "9999"
	}

	gou.Debug("listening at " + port)
	http.ListenAndServe(":"+port, nil)
}

func RecoverAndLogError() {
	if err := recover(); err != nil {
		out := []byte{}
		runtime.Stack(out, true)
		gou.Error(string(out))
	}
}
