package scoring

import (
	"testing"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

func getTestIssueScorer(closedOpenRatio float64, closedIssueCount int, openIssueCount int) IssueScorer {
	issues := make([]*github.Issue, openIssueCount+closedIssueCount)
	closed := "closed"
	open := "open"
	for i := 0; i < closedIssueCount; i++ {
		issues = append(issues, &github.Issue{State: &closed})
	}
	for i := 0; i < openIssueCount; i++ {
		issues = append(issues, &github.Issue{State: &open})
	}

	return IssueScorer{
		data: helpers.GithubRepoInfo{
			Issues: issues,
		},
		config: config.IssuesScoringConfig{
			MaxPenalty:      2.0,
			ClosedOpenRatio: closedOpenRatio,
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
