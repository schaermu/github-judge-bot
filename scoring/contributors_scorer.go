package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

// ContributorsScorer provides a scoring based on contributors on a repository.
type ContributorsScorer struct {
	Scorer
	data   *data.GithubRepoInfo
	config config.ScorerConfig
}

// GetScore calculates a score based on the active contributors.
// For each missing contributor compared to min_contributors, a percentage of the configured max_penalty is applied to the score.
func (s ContributorsScorer) GetScore(currentScore float64, penalties []Penalty) (float64, []Penalty) {
	// we calculate the percentage of contributors vs. required contributors and apply that percentage as a penalty
	minContribs := s.config.GetInt("min_contributors")
	percentage := 100 / float64(minContribs) * float64(len(s.data.Contributors))
	scoreChange := 0.0
	if percentage < 100 {
		unrounded := (100 - percentage) * (s.config.MaxPenalty / 100)
		// in order to prevent weird score changes like -1.98, we ceil to the first decimal
		scoreChange = math.Ceil(unrounded*10) / 10
	}

	if scoreChange > 0 {
		currentScore -= scoreChange

		penalties = append(penalties, Penalty{
			ScorerName: "Contributors",
			Reason:     fmt.Sprintf("There are only %d out of %d required contributors", len(s.data.Contributors), minContribs),
			Amount:     scoreChange,
		})
	}

	return currentScore, penalties
}
