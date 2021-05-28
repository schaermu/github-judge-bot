package scoring

import (
	"fmt"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

// StarsScorer provides a scoring based on the repository's star gazers.
type StarsScorer struct {
	data   *data.GithubRepoInfo
	config config.ScorerConfig
}

// GetScore calculates a score based on star gazers. If a repository is below the threshold configured in "min_stars", a penalty is applied.
func (s StarsScorer) GetScore(currentScore float64, penalties []Penalty) (float64, []Penalty) {
	reqStars := s.config.GetInt("min_stars")
	if s.data.Stars < reqStars {
		scoreChange := s.config.MaxPenalty
		currentScore -= scoreChange

		penalties = append(penalties, Penalty{
			ScorerName: "Stars",
			Reason:     fmt.Sprintf("Less than %d stars (%d stars)", reqStars, s.data.Stars),
			Amount:     scoreChange,
		})
	}

	return currentScore, penalties
}
