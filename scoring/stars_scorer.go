package scoring

import (
	"fmt"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

type StarsScorer struct {
	data   data.GithubRepoInfo
	config config.ScorerConfig
}

func (s StarsScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
	reqStars := s.config.GetInt("min_stars")
	if s.data.Stars < reqStars {
		scoreChange := s.config.MaxPenalty
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("Less than %d stars (*%d* stars)", reqStars, s.data.Stars),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}
