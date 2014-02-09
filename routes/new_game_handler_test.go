package routes_test

import (
	. "github.com/artemave/conways-go/dependencies/ginkgo"
	. "github.com/artemave/conways-go/dependencies/gomega"
	. "github.com/artemave/conways-go/routes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("NewGameHandler", func() {
	server := httptest.NewServer(nil)
	RegisterRoutes()

	Describe("StartNewGameHandler", func() {
		It("redirects from root to new game url", func() {
			resp, err := http.Get(server.URL)
			if err != nil {
				Fail(err.Error())
			}
			Expect(resp.Request.URL.String()).To(MatchRegexp("/games/.+"))
		})
	})

	It("returns html/javascript for the browser to kick off websocket session", func() {
		resp, err := http.Get(server.URL + "/games/234234asd")
		if err != nil {
			Fail(err.Error())
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Fail(err.Error())
		}
		resp.Body.Close()
		Expect(string(body)).To(MatchRegexp("public/bundle.js"))
	})
})
