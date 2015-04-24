package scores

import (
	"errors"
	"sync"

	gga "github.com/artemave/conways-go/google_games_adapter"
)

// FetchLeaderboards : from Google Games API
func FetchLeaderboards(gapi interface {
	Leaderboards() ([]gga.Leaderboard, error)
	Scores(gga.Leaderboard) (*gga.LeaderboardScores, error)
}) ([]*gga.LeaderboardScores, error) {

	leaderboards, err := gapi.Leaderboards()
	if err != nil {
		return nil, err
	}

	var boardsScores []*gga.LeaderboardScores

	var wg sync.WaitGroup
	res := make(chan *gga.LeaderboardScores, len(leaderboards))
	errs := make(chan error, len(leaderboards))

	defer close(errs)
	defer close(res)

	wg.Add(len(leaderboards))

	for _, board := range leaderboards {
		go func(b gga.Leaderboard) {
			defer wg.Done()

			scores, err := gapi.Scores(b)
			if err != nil {
				errs <- err
			}
			res <- scores
		}(board)
	}
	wg.Wait()

FORZ:
	for {
		select {
		case err := <-errs:
			return nil, err
		case scores := <-res:
			boardsScores = append(boardsScores, scores)
		default:
			break FORZ
		}
	}

	return boardsScores, nil
}

type game interface {
	Size() string
	SetScoredBy(string)
	IsPractice() bool
	IsFinished() bool
	GetScoredBy() *string
}

// SubmitScore : to Google Games API
func SubmitScore(gapi interface {
	Leaderboards() ([]gga.Leaderboard, error)
	CurrentPlayerScore(gga.Leaderboard) (*gga.PlayerScore, error)
	SubmitScore(gga.Leaderboard, int64) error
}, g game) error {
	err := validateGameSubmitScore(g)
	if err != nil {
		return err
	}

	leaderboards, err := gapi.Leaderboards()
	if err != nil {
		return err
	}

	for _, board := range leaderboards {
		if board.Name == g.Size() {
			score, err := gapi.CurrentPlayerScore(board)
			if err != nil {
				return err
			}

			err = gapi.SubmitScore(board, score.ScoreValue+3)
			if err != nil {
				return err
			}
			g.SetScoredBy(score.PlayerId)
			break
		}
	}

	return nil
}

func validateGameSubmitScore(game game) error {
	if game.IsPractice() {
		return errors.New("Score can not be submitted for practice game.")
	}

	if !game.IsFinished() {
		return errors.New("Can't submit score for a game that has not yet finished.")
	}

	if game.GetScoredBy() != nil {
		return errors.New("Score has already been submitted.")
	}

	return nil
}
