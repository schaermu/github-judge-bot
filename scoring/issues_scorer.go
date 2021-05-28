package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

// IssuesScorer provides a scoring based on the open and closed issues in a repository.
type IssuesScorer struct {
	data   *data.GithubRepoInfo
	config config.ScorerConfig
}

// GetScore calculates a score based on issues statistics.
// The scoring is determined based on the ratio between open and closed issues (basically how many open issues are allowed per closed ones).
// If this ratio is above the threshold specified in "closed_open_ratio", a penalty is applied.
func (s IssuesScorer) GetScore(currentScore float64, penalties []Penalty) (float64, []Penalty) {
	open := 0.0
	closed := 0.0
	for _, issue := range s.data.Issues {
		if issue.GetState() == "open" {
			open++
		}
		if issue.GetState() == "closed" {
			closed++
		}
	}

	ratio := open / closed
	scoreChange := s.config.MaxPenalty
	requiredRatio := s.config.GetFloat64("closed_open_ratio")
	if ratio <= requiredRatio {
		scoreChange = 0
	}

	if !math.IsNaN(ratio) && scoreChange > 0 {
		currentScore -= scoreChange

		penalties = append(penalties, Penalty{
			ScorerName: "Issues",
			Reason:     fmt.Sprintf("Closed-Open Ratio on issues is above 1:%.2f (*1:%.2f*)", requiredRatio, ratio),
			Amount:     scoreChange,
		})
	}

	return currentScore, penalties
}
