package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

type IssueScorer struct {
	data   helpers.GithubRepoInfo
	config config.IssuesScoringConfig
}

func (s IssueScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
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

	ratio := (closed / open) * 100
	scoreChange := s.config.MaxPenalty
	if ratio > s.config.OpenClosedRatio {
		scoreChange = 0
	}

	if !math.IsNaN(ratio) && scoreChange > 0 {
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("Closed-Open Ratio on issues is below %.2f%% (*%.2f%%*)", s.config.OpenClosedRatio, ratio),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}
