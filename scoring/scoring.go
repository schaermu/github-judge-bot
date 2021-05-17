package scoring

import (
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

type ScoringPenalty struct {
	Reason string
	Amount float64
}

type Scorer interface {
	GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty)
}

func CreateScorer(data data.GithubRepoInfo, config config.ScorerConfig) Scorer {
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

func CreateScorerMap(data data.GithubRepoInfo, configs []config.ScorerConfig) (scorers map[string]Scorer, score float64) {
	// create map of all scorers and initialize score to maximum possible
	scorers = map[string]Scorer{}
	for _, config := range configs {
		scorers[config.Name] = CreateScorer(data, config)
		score += config.MaxPenalty
	}
	return
}

func GetTotalScore(data data.GithubRepoInfo, scorers []config.ScorerConfig) (score float64, maxScore float64, penalties []ScoringPenalty) {
	scorerMap, maxScore := CreateScorerMap(data, scorers)
	score = maxScore
	// execute scorers
	for _, scorer := range scorerMap {
		score, penalties = scorer.GetScore(score, penalties)
	}

	return score, maxScore, penalties
}
