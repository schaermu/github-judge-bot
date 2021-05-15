package scoring

import (
	"testing"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

func getTestScorer(stars int, minStars int) StarsScorer {
	return StarsScorer{
		data: helpers.GithubRepoInfo{
			Stars: stars,
		},
		config: config.StarsScoringConfig{
			MaxPenalty: 2.0,
			MinStars:   minStars,
		},
	}
}

func TestGetScore(t *testing.T) {
	scorer := getTestScorer(800, 600)

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

func TestGetScorePenalty(t *testing.T) {
	scorer := getTestScorer(800, 801)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 8.0
	expectedPenalty := 2.0

	if score == 10 {
		t.Errorf("GetScore() failed, expected score to be %.2f, got %.2f", expectedScore, score)
	}
	if len(penalties) == 0 {
		t.Errorf("GetScore() failed, expected 1 penalty, got %d", len(penalties))
	}
	if penalties[0].Amount != expectedPenalty {
		t.Errorf("GetScore() failed, expected penalty amount of %.2f, got %.2f", expectedPenalty, penalties[0].Amount)
	}
}
