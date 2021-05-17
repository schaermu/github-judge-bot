package scoring

import (
	"fmt"
	"testing"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/stretchr/testify/assert"
)

func getTestStarsScorer(stars int, minStars int) StarsScorer {
	return StarsScorer{
		data: data.GithubRepoInfo{
			Stars: stars,
		},
		config: config.ScorerConfig{
			MaxPenalty: 2.0,
			Settings: map[string]string{
				"min_stars": fmt.Sprintf("%d", minStars),
			},
		},
	}
}

func TestStarsScorerGetScore(t *testing.T) {
	scorer := getTestStarsScorer(800, 600)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 10.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 0)
}

func TestStarsScorerGetScorePenalty(t *testing.T) {
	scorer := getTestStarsScorer(800, 900)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 8.0
	expectedPenalty := 2.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 1)
	assert.Equal(t, expectedPenalty, penalties[0].Amount)
}
