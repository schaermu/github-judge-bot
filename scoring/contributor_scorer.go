package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

type ContributorScorer struct {
	data   helpers.GithubRepoInfo
	config config.ContributorsConfig
}

func (s ContributorScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
	// we calculate the percentage of contributors vs. required contributors and apply that percentage as a penalty
	percentage := 100 / float64(s.config.MinContributors) * float64(len(s.data.Contributors))
	scoreChange := 0.0
	if percentage < 100 {
		unrounded := (100 - percentage) * (s.config.MaxPenalty / 100)
		// in order to prevent weird score changes like -1.98, we ceil to the first decimal
		scoreChange = math.Ceil(unrounded*10) / 10
	}

	if scoreChange > 0 {
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("There are only *%d/%d* required contributors", len(s.data.Contributors), s.config.MinContributors),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}