package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

type IssuesScorer struct {
	data   data.GithubRepoInfo
	config config.ScorerConfig
}

func (s IssuesScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
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

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("Closed-Open Ratio on issues is below 1:%.2f (*1:%.2f*)", requiredRatio, ratio),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}
