package scoring

import (
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

// Summary contains the result of a full scoring run.
type Summary struct {
	Score          float64
	MaxScore       float64
	TotalPenalties float64
	Penalties      []Penalty
}

// Penalty contains details about a penalty that was applied.
type Penalty struct {
	ScorerName string
	Reason     string
	Amount     float64
}

// Scorer provides the interface all scorers must follow.
type Scorer interface {
	GetScore(currentScore float64, penalties []Penalty) (float64, []Penalty)
}

// CreateScorer builds a scorer object for a specific name configured in the config.
func CreateScorer(data *data.GithubRepoInfo, config config.ScorerConfig) Scorer {
	switch config.Name {
	case "stars":
		return StarsScorer{data: data, config: config}
	case "issues":
		return IssuesScorer{data: data, config: config}
	case "commit-activity":
		return CommitActivityScorer{data: data, config: config}
	case "contributors":
		return ContributorsScorer{data: data, config: config}
	case "license":
		return LicenseScorer{data: data, config: config}
	default:
		return nil
	}
}

// CreateScorerMap builds a map of scorer names/objects and calculates the maximum score.
func CreateScorerMap(data *data.GithubRepoInfo, configs []config.ScorerConfig) (scorers map[string]Scorer, score float64) {
	// create map of all scorers and initialize score to maximum possible
	scorers = map[string]Scorer{}
	for _, config := range configs {
		scorers[config.Name] = CreateScorer(data, config)
		score += config.MaxPenalty
	}
	return
}

// GetTotalScore calculates the score summary for a github repository with the specified scorers.
func GetTotalScore(data *data.GithubRepoInfo, scorers []config.ScorerConfig) (summary Summary) {
	scorerMap, maxScore := CreateScorerMap(data, scorers)
	score := maxScore
	penalties := []Penalty{}
	// run scorers
	for _, scorer := range scorerMap {
		score, penalties = scorer.GetScore(score, penalties)
	}

	return Summary{
		Score:          score,
		MaxScore:       maxScore,
		TotalPenalties: maxScore - score,
		Penalties:      penalties,
	}
}
