package scoring

import (
	"math/big"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

type Scorer interface {
	Score(currentScore *big.Rat, penalties []ScoringPenalty, data helpers.GithubRepoInfo)
}

func GetTotalScore(data helpers.GithubRepoInfo, scoreConfig config.ScoringConfig) (score float64, penalties []ScoringPenalty) {
	stars := StarsScorer{data: data, config: scoreConfig.Stars}
	score, penalties = stars.GetScore(score, penalties)

	return score, penalties
}
