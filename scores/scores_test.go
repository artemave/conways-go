package scores_test

import (
	gga "github.com/artemave/conways-go/google_games_adapter"
	"github.com/artemave/conways-go/scores"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type submitScoreCall struct {
	Board gga.Leaderboard
	Score int64
}

type gapiSpy struct {
	SubmitScoreCalls []submitScoreCall
	PlayerScore      *gga.PlayerScore
}

func (gapi *gapiSpy) Leaderboards() ([]gga.Leaderboard, error) {
	leaderbords := []gga.Leaderboard{
		gga.Leaderboard{Name: "small"},
		gga.Leaderboard{Name: "medium"},
		gga.Leaderboard{Name: "large"},
	}
	return leaderbords, nil
}

func (gapi *gapiSpy) CurrentPlayerScore(gga.Leaderboard) (*gga.PlayerScore, error) {
	return gapi.PlayerScore, nil
}

func (gapi *gapiSpy) SubmitScore(board gga.Leaderboard, score int64) error {
	gapi.SubmitScoreCalls = append(gapi.SubmitScoreCalls, submitScoreCall{board, score})
	return nil
}

type gameSpy struct {
	scoredBy   string
	isPractice bool
	isFinished bool
	size       string
}

func (g *gameSpy) Size() string {
	return g.size
}

func (g *gameSpy) SetSize(size string) {
	g.size = size
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

func (g *gameSpy) GetScoredBy() string {
	return g.scoredBy
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
			game = &gameSpy{isPractice: false, isFinished: true, size: "medium"}
			gapi = &gapiSpy{
				PlayerScore: &gga.PlayerScore{
					PlayerId:   "artem",
					ScoreValue: 3,
				},
			}
		})
		Describe("validations", func() {
			validationAssertions := func() {
				It("returns an error", func() {
					err := scores.SubmitScore(gapi, game)
					Expect(err).ToNot(BeNil())
				})
				It("doesn't submit the score", func() {
					_ = scores.SubmitScore(gapi, game)
					Expect(len(gapi.SubmitScoreCalls)).To(Equal(0))
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

		It("submits without errors and only once", func() {
			err := scores.SubmitScore(gapi, game)

			Expect(err).To(BeNil())
			Expect(len(gapi.SubmitScoreCalls)).To(Equal(1))
		})

		It("submits to a board corresponding to game size", func() {
			game.SetSize("small")

			_ = scores.SubmitScore(gapi, game)

			Expect(gapi.SubmitScoreCalls[0].Board.Name).To(Equal("small"))
		})

		It("increments current player score by 3", func() {
			gapi.PlayerScore.ScoreValue = 6

			_ = scores.SubmitScore(gapi, game)

			Expect(gapi.SubmitScoreCalls[0].Score).To(Equal(int64(9)))
		})
	})
})
