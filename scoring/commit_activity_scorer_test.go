package scoring

import (
	"testing"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

func getTestCommitActivityScorer(inactiveWeekCount int, penaltyPerWeek float64) CommitActivityScorer {
	// prepare commit activity data
	activity := make([]*github.WeeklyCommitActivity, 0)
	if inactiveWeekCount > 0 {
		zero := 0
		for i := 1; i < inactiveWeekCount+1; i++ {
			time := github.Timestamp{Time: time.Now().Local().AddDate(0, 0, -7*i)}
			activity = append(activity, &github.WeeklyCommitActivity{Week: &time, Total: &zero})
		}
	}
	one := 1
	firstActiveTime := github.Timestamp{Time: time.Now().Local().AddDate(0, 0, -7*inactiveWeekCount+2)}
	activity = append(activity, &github.WeeklyCommitActivity{Week: &firstActiveTime, Total: &one})

	return CommitActivityScorer{
		data: helpers.GithubRepoInfo{
			CommitActivity: activity,
		},
		config: config.CommitActivityConfig{
			MaxPenalty:              2.0,
			WeeklyInactivityPenalty: penaltyPerWeek,
		},
	}
}

func TestActivityScorerGetScore(t *testing.T) {
	scorer := getTestCommitActivityScorer(0, 0.25)

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

func TestActivityScorerGetScorePenalty(t *testing.T) {
	scorer := getTestCommitActivityScorer(2, 0.25)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 9.5
	expectedPenalty := 0.5

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

func TestActivityScorerGetScoreCappedPenalty(t *testing.T) {
	scorer := getTestCommitActivityScorer(40, 0.25)

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
