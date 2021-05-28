package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

// CommitActivityScorer provides a scoring based on commits made within the last year.
type CommitActivityScorer struct {
	data   *data.GithubRepoInfo
	config config.ScorerConfig
}

// GetScore calculates a score based on the commit activity.
// For each week of inactivity, the penalty configured in the setting "weekly_penalty" is applied until a maximum of max_penalty.
func (s CommitActivityScorer) GetScore(currentScore float64, penalties []Penalty) (float64, []Penalty) {
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

		penalties = append(penalties, Penalty{
			ScorerName: "CommitActivity",
			// TODO: pluralization
			Reason: fmt.Sprintf("The last commit was more than %d week(s) ago", weeksWithoutActivity),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}
