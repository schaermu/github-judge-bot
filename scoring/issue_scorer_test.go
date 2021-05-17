package scoring

import (
	"fmt"
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
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
		data: helpers.GithubRepoInfo{
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

func TestIssueScorerGetScorePenalty(t *testing.T) {
	scorer := getTestIssueScorer(0.2, 20, 10)

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
