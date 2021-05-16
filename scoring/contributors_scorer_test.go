package scoring

import (
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

func getTestContributorScorer(contributorCount int, minContributors int) ContributorScorer {
	// prepare commit activity data
	contributors := make([]*github.ContributorStats, 0)
	if contributorCount > 0 {
		for i := 1; i < contributorCount+1; i++ {
			contributors = append(contributors, &github.ContributorStats{})
		}
	}

	return ContributorScorer{
		data: helpers.GithubRepoInfo{
			Contributors: contributors,
		},
		config: config.ContributorsConfig{
			MaxPenalty:      2.0,
			MinContributors: minContributors,
		},
	}
}

func TestContributorScorerGetScore(t *testing.T) {
	scorer := getTestContributorScorer(5, 3)

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

func TestContributorScorerGetScorePenalty(t *testing.T) {
	scorer := getTestContributorScorer(6, 10)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 9.2
	expectedPenalty := 0.8

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

func TestContributorScorerGetScoreCappedPenalty(t *testing.T) {
	scorer := getTestContributorScorer(1, 100)

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
