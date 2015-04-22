package googleGamesAdapter

import (
	"net/http"

	gg "google.golang.org/api/games/v1"
)

// Games scope url
var GamesScope = gg.GamesScope

// GoogleGamesAPIAdapter : google games api wrapper
type GoogleGamesAPIAdapter struct {
	gapi *gg.Service
}

// New : new google games api client
func New(client *http.Client) (*GoogleGamesAPIAdapter, error) {
	gapi, err := gg.New(client)
	if err != nil {
		return nil, err
	}
	gga := &GoogleGamesAPIAdapter{gapi}
	return gga, nil
}

// Leaderboard resource
type Leaderboard struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// LeaderboardScores resource
type LeaderboardScores struct {
	Leaderboard
	PlayerScore *gg.LeaderboardEntry   `json:"playerScore"`
	Items       []*gg.LeaderboardEntry `json:"items"`
}

// Leaderboards : fetch leaderboards from google games api
func (gga *GoogleGamesAPIAdapter) Leaderboards() ([]Leaderboard, error) {
	gboards, err := gga.gapi.Leaderboards.List().Do()
	if err != nil {
		return nil, err
	}

	var leaderboards []Leaderboard
	for _, board := range gboards.Items {
		b := Leaderboard{
			Id:   board.Id,
			Name: board.Name,
		}
		leaderboards = append(leaderboards, b)
	}

	return leaderboards, nil
}

// Scores : fetch board scores from google games api
func (gga *GoogleGamesAPIAdapter) Scores(board Leaderboard) (*LeaderboardScores, error) {
	scores := &LeaderboardScores{Leaderboard: board}
	s, err := gga.gapi.Scores.List(board.Id, "PUBLIC", "ALL_TIME").MaxResults(10).Do()
	if err != nil {
		return nil, err
	}
	scores.PlayerScore = s.PlayerScore
	scores.Items = s.Items
	return scores, nil
}

type PlayerScore struct {
	PlayerId   string `json:"playerId"`
	ScoreValue int64  `json:"scoreValue"`
}

//
func (gga *GoogleGamesAPIAdapter) CurrentPlayerScore(board Leaderboard) (*PlayerScore, error) {
	s, err := gga.gapi.Scores.Get("me", board.Id, "ALL_TIME").Do()
	if err != nil {
		return nil, err
	}

	var currentScore int64
	if len(s.Items) == 0 {
		currentScore = 0
	} else {
		currentScore = int64(s.Items[0].ScoreValue)
	}

	playerScore := &PlayerScore{
		PlayerId:   s.Player.PlayerId,
		ScoreValue: currentScore,
	}
	return playerScore, nil
}

//
func (gga *GoogleGamesAPIAdapter) SubmitScore(board Leaderboard, newScoreValue int64) error {
	_, err := gga.gapi.Scores.Submit(board.Id, newScoreValue).Do()
	if err != nil {
		return err
	}
	return nil
}
