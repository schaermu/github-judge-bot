package scoring

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/stretchr/testify/assert"
)

func getTestCommitActivityData(inactiveWeekCount int) []*github.WeeklyCommitActivity {
	// prepare commit activity data: first an active week, then append inactive ones (we loop in reverse)
	activity := make([]*github.WeeklyCommitActivity, 0)
	one := 1
	firstActiveTime := github.Timestamp{Time: time.Now().Local().AddDate(0, 0, -7*inactiveWeekCount+2)}
	activity = append(activity, &github.WeeklyCommitActivity{Week: &firstActiveTime, Total: &one})
	if inactiveWeekCount > 0 {
		zero := 0
		for i := 1; i < inactiveWeekCount+1; i++ {
			time := github.Timestamp{Time: time.Now().Local().AddDate(0, 0, -7*i)}
			activity = append(activity, &github.WeeklyCommitActivity{Week: &time, Total: &zero})
		}
	}
	return activity
}

func getTestCommitActivityScorer(inactiveWeekCount int, penaltyPerWeek float64) CommitActivityScorer {
	return CommitActivityScorer{
		data: data.GithubRepoInfo{
			CommitActivity: getTestCommitActivityData(inactiveWeekCount),
		},
		config: config.ScorerConfig{
			MaxPenalty: 2.0,
			Settings: map[string]string{
				"weekly_penalty": fmt.Sprintf("%.2f", penaltyPerWeek),
			},
		},
	}
}

func TestCommitActivityScorerGetScore(t *testing.T) {
	scorer := getTestCommitActivityScorer(0, 0.25)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 10.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 0)
}

func TestCommitActivityScorerGetScorePenalty(t *testing.T) {
	scorer := getTestCommitActivityScorer(2, 0.25)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 9.5
	expectedPenalty := 0.5

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 1)
	assert.Equal(t, expectedPenalty, penalties[0].Amount)
}

func TestCommitActivityScorerGetScoreCappedPenalty(t *testing.T) {
	scorer := getTestCommitActivityScorer(40, 0.25)

	penalties := make([]ScoringPenalty, 0)
	score := 10.0
	score, penalties = scorer.GetScore(score, penalties)

	expectedScore := 8.0
	expectedPenalty := 2.0

	assert.Equal(t, expectedScore, score)
	assert.Len(t, penalties, 1)
	assert.Equal(t, expectedPenalty, penalties[0].Amount)
}
