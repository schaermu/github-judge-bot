package scoring

import (
	"math/big"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

type ScoringPenalty struct {
	Reason string
	Amount float64
}

type Scorer interface {
	Score(currentScore *big.Rat, penalties []ScoringPenalty, data helpers.GithubRepoInfo)
}

func GetTotalScore(data helpers.GithubRepoInfo, scoreConfig config.ScoringConfig) (score float64, penalties []ScoringPenalty) {
	score = scoreConfig.MaxScore

	stars := StarsScorer{data: data, config: scoreConfig.Stars}
	score, penalties = stars.GetScore(score, penalties)

	issues := IssueScorer{data: data, config: scoreConfig.Issues}
	score, penalties = issues.GetScore(score, penalties)

	commitActivity := CommitActivityScorer{data: data, config: scoreConfig.CommitActivity}
	score, penalties = commitActivity.GetScore(score, penalties)

	contributors := ContributorScorer{data: data, config: scoreConfig.Contributors}
	score, penalties = contributors.GetScore(score, penalties)

	license := LicenseScorer{data: data, config: scoreConfig.License}
	score, penalties = license.GetScore(score, penalties)

	return score, penalties
}
