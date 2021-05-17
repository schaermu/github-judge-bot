package scoring

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
	"github.com/stretchr/testify/assert"
)

func getTestContributorData(contributorCount int) []*github.ContributorStats {
	// prepare commit activity data
	contributors := make([]*github.ContributorStats, 0)
	if contributorCount > 0 {
		for i := 1; i < contributorCount+1; i++ {
			contributors = append(contributors, &github.ContributorStats{})
		}
	}
	return contributors
}

func getTestContributorScorer(contributorCount int, minContributors int) ContributorsScorer {
	return ContributorsScorer{
		data: helpers.GithubRepoInfo{
			Contributors: getTestContributorData(contributorCount),
		},
		config: config.ScorerConfig{
			MaxPenalty: 2.0,
			Settings: map[string]string{
				"min_contributors": fmt.Sprintf("%d", minContributors),
			},
		},
	}
}

func TestContributorScorerGetScore(t *testing.T) {
	scorer := getTestContributorScorer(5, 3)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 10.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 0)
}

func TestContributorScorerGetScorePenalty(t *testing.T) {
	scorer := getTestContributorScorer(6, 10)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 9.2
	expectedPenalty := 0.8

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 1)
	assert.Equal(t, expectedPenalty, penalties[0].Amount)
}

func TestContributorScorerGetScoreCappedPenalty(t *testing.T) {
	scorer := getTestContributorScorer(1, 100)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 8.0
	expectedPenalty := 2.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 1)
	assert.Equal(t, expectedPenalty, penalties[0].Amount)
}
