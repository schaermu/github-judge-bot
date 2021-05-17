package scoring

import (
	"strings"
	"testing"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

func getTestLicenseScorer(license string, validLicenses []string) LicenseScorer {
	return LicenseScorer{
		data: helpers.GithubRepoInfo{
			License:   license,
			LicenseId: license,
		},
		config: config.ScorerConfig{
			MaxPenalty: 2.0,
			Settings: map[string]string{
				"valid_license_ids": strings.Join(validLicenses, ","),
			},
		},
	}
}

func TestLicenseScorerGetScore(t *testing.T) {
	scorer := getTestLicenseScorer("MIT", nil)

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

func TestLicenseScorerGetScorePenalty(t *testing.T) {
	scorer := getTestLicenseScorer("MIT", []string{"BSD", "Apache-2.0"})

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
