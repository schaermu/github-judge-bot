package scoring

import (
	"fmt"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

type StarsScorer struct {
	data   helpers.GithubRepoInfo
	config config.StarsScoringConfig
}

func (s StarsScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
	if s.data.Stars < s.config.MinStars {
		scoreChange := s.config.MaxPenalty
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("Less than %d stars (*%d* stars)", s.config.MinStars, s.data.Stars),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}
