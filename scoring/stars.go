package scoring

import (
	"fmt"
	"math"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

type StarsScorer struct {
	data   helpers.GithubRepoInfo
	config config.StarsScoringConfig
}

func (s StarsScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
	starRel := float64(s.data.Stars/s.config.MinStars) * s.config.MaxPenalty
	if starRel > s.config.MaxPenalty {
		starRel = s.config.MaxPenalty
	}

	if math.Abs(starRel-s.config.MaxPenalty) > 0 {
		scoreChange := math.Abs(starRel - s.config.MaxPenalty)
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("Less than %d stars (*%d* stars)", s.config.MinStars, s.data.Stars),
			Amount: currentScore,
		})
	}

	return currentScore, penalties
}
