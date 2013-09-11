package main

import (
  "github.com/artemave/conways-go/game"
  "github.com/araddon/gou"
  "net/http"
  "io/ioutil"
  "os"
  "log"
  "encoding/json"
)

func main() {
  g := game.Game{Rows: 500, Cols: 500}

  port := os.Getenv("PORT")
  if (port == "") {
    port = "9999"
  }

  gou.SetLogger(log.New(os.Stderr, "", log.LstdFlags), "debug")

  http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    gou.Debug("POST /go")

    // Get json from request
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      gou.Error("error: ", err)
    } else {
      gou.Debug("request body: ", string(body))
    }

    // Unmarshal to points
    points := []game.Point{}
    err = json.Unmarshal(body, &points)
    if err != nil {
      gou.Error("error: ", err)
    }

    // Create generation from points
    current_generation := g.PointsToGeneration(&points)

    // Calculate next generation
    next_generation := g.NextGeneration(current_generation)

    // Dump it to json
    b, err := json.Marshal(next_generation)
    if err != nil {
      gou.Error("error: ", err)
    }

    w.Write(b)
  })

  http.Handle("/", http.FileServer(http.Dir("./public/")))

  gou.Debug("listening at " + port)
  http.ListenAndServe(":" + port, nil)
}
