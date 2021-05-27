package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

type CommitActivityScorer struct {
	data   *data.GithubRepoInfo
	config config.ScorerConfig
}

func (s CommitActivityScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
	weeksWithoutActivity := 0
	// loop in reverse because we get oldest first from github
	for i := len(s.data.CommitActivity) - 1; i >= 0; i-- {
		commit := s.data.CommitActivity[i]
		if commit.GetTotal() == 0 {
			weeksWithoutActivity++
		} else {
			break
		}
	}

	if weeksWithoutActivity > 0 {
		scoreChange := math.Min(float64(weeksWithoutActivity)*s.config.GetFloat64("weekly_penalty"), s.config.MaxPenalty)
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("The last commit was more than *%d* week(s) ago", weeksWithoutActivity),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}
