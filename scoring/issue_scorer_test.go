package scoring

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/stretchr/testify/assert"
)

func getTestIssueData(closedOpenRatio float64, closedIssueCount int, openIssueCount int) []*github.Issue {
	issues := make([]*github.Issue, openIssueCount+closedIssueCount)
	closed := "closed"
	open := "open"
	for i := 0; i < closedIssueCount; i++ {
		issues = append(issues, &github.Issue{State: &closed})
	}
	for i := 0; i < openIssueCount; i++ {
		issues = append(issues, &github.Issue{State: &open})
	}
	return issues
}

func getTestIssueScorer(closedOpenRatio float64, closedIssueCount int, openIssueCount int) IssuesScorer {
	return IssuesScorer{
		data: &data.GithubRepoInfo{
			Issues: getTestIssueData(closedOpenRatio, closedIssueCount, openIssueCount),
		},
		config: config.ScorerConfig{
			MaxPenalty: 2.0,
			Settings: map[string]string{
				"closed_open_ratio": fmt.Sprintf("%.2f", closedOpenRatio),
			},
		},
	}
}

func TestIssueScorerGetScore(t *testing.T) {
	scorer := getTestIssueScorer(0.2, 20, 4)

	penalties := make([]Penalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 10.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 0)
}

func TestIssueScorerGetScorePenalty(t *testing.T) {
	scorer := getTestIssueScorer(0.2, 20, 10)

	penalties := make([]Penalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 8.0
	expectedPenalty := 2.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 1)
	assert.Equal(t, expectedPenalty, penalties[0].Amount)
}
