package scores_test

import (
	gga "github.com/artemave/conways-go/google_games_adapter"
	"github.com/artemave/conways-go/scores"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type gapiSpy struct {
	SubmitScoreCallNumber int
}

func (gapi *gapiSpy) Leaderboards() ([]gga.Leaderboard, error) {
	return []gga.Leaderboard{}, nil
}

func (gapi *gapiSpy) CurrentPlayerScore(gga.Leaderboard) (*gga.PlayerScore, error) {
	return nil, nil
}

func (gapi *gapiSpy) SubmitScore(gga.Leaderboard, int64) error {
	gapi.SubmitScoreCallNumber++
	return nil
}

type gameSpy struct {
	scoredBy   string
	isPractice bool
	isFinished bool
}

func (g *gameSpy) Size() string {
	return "small"
}

func (g *gameSpy) SetScoredBy(player string) {
	g.scoredBy = player
}

func (g *gameSpy) IsPractice() bool {
	return g.isPractice
}

func (g *gameSpy) IsFinished() bool {
	return g.isFinished
}

func (g *gameSpy) GetScoredBy() *string {
	return &g.scoredBy
}

func (g *gameSpy) SetIsPractice(v bool) {
	g.isPractice = v
}

func (g *gameSpy) SetIsFinished(v bool) {
	g.isFinished = v
}

var _ = Describe("Scores", func() {
	var game *gameSpy
	var gapi *gapiSpy

	Describe("SubmitScore", func() {
		BeforeEach(func() {
			game = &gameSpy{isPractice: false, isFinished: true}
			gapi = &gapiSpy{}
		})
		Describe("validations", func() {
			validationAssertions := func() {
				It("returns an error", func() {
					err := scores.SubmitScore(gapi, game)
					Expect(err).ToNot(BeNil())
				})
				It("doesn't submit the score", func() {
					_ = scores.SubmitScore(gapi, game)
					Expect(gapi.SubmitScoreCallNumber).To(Equal(0))
				})
			}
			Context("Score has already been submitted for this game", func() {
				BeforeEach(func() {
					game.SetScoredBy("player")
				})
				validationAssertions()
			})
			Context("Practice game", func() {
				BeforeEach(func() {
					game.SetIsPractice(true)
				})
				validationAssertions()
			})
			Context("Game is not finished yet", func() {
				BeforeEach(func() {
					game.SetIsFinished(false)
				})
				validationAssertions()
			})
		})
	})
})
