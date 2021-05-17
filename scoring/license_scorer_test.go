package scoring

import (
	"strings"
	"testing"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/stretchr/testify/assert"
)

func getTestLicenseScorer(license string, validLicenses []string) LicenseScorer {
	return LicenseScorer{
		data: data.GithubRepoInfo{
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

	expectedScore := 10.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 0)
}

func TestLicenseScorerGetScorePenalty(t *testing.T) {
	scorer := getTestLicenseScorer("MIT", []string{"BSD", "Apache-2.0"})

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 8.0
	expectedPenalty := 2.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 1)
	assert.Equal(t, expectedPenalty, penalties[0].Amount)
}
