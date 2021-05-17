package scoring

import (
	"fmt"
	"testing"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

func getTestStarsScorer(stars int, minStars int) StarsScorer {
	return StarsScorer{
		data: helpers.GithubRepoInfo{
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
	if score != 10 {
		t.Errorf("GetScore() failed, expected score to be 10, got %.2f", score)
	}
	if len(penalties) > 0 {
		t.Errorf("GetScore() failed, expected no penalties, got %d", len(penalties))
	}
}

func TestStarsScorerGetScorePenalty(t *testing.T) {
	scorer := getTestStarsScorer(800, 900)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 8.0
	expectedPenalty := 2.0

	if score != expectedScore {
		t.Errorf("GetScore() failed, expected score to be %.2f, got %.2f", expectedScore, score)
	}
	if len(penalties) == 0 {
		t.Errorf("GetScore() failed, expected 1 penalty, got %d", len(penalties))
	}
	if penalties[0].Amount != expectedPenalty {
		t.Errorf("GetScore() failed, expected penalty amount of %.2f, got %.2f", expectedPenalty, penalties[0].Amount)
	}
}
