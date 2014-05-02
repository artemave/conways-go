package routes_test

import (
	"net/http"
	"net/http/httptest"
	. "github.com/artemave/conways-go/dependencies/ginkgo"
	. "github.com/artemave/conways-go/dependencies/gomega"
)

var _ = Describe("NewGameHandler", func() {
	server := httptest.NewServer(nil)

	Describe("StartNewGameHandler", func() {
		It("redirects from root to new game url", func() {
			resp, err := http.Get(server.URL)
			if err != nil {
				Fail(err.Error())
			}
			defer resp.Body.Close()
			Expect(resp.Request.URL.String()).To(MatchRegexp("/games/.+"))
		})
	})

	// commented out because it serves static file and server can only find it if run from project root
	//
	// It("returns html/javascript for the browser to kick off websocket session", func() {
	// 	resp, err := http.Get(server.URL)
	// 	if err != nil {
	// 		Fail(err.Error())
	// 	}
	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		Fail(err.Error())
	// 	}
	// 	defer resp.Body.Close()
	// 	Expect(string(body)).To(MatchRegexp("public/bundle.js"))
	// })
})
