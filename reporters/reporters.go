package reporters

import (
	"regexp"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/schaermu/go-github-judge-bot/helpers"
	"github.com/schaermu/go-github-judge-bot/scoring"
)

var GITHUB_URL_MATCHER = regexp.MustCompile("github.com/([^/]+)/([^/<>]+)")

type Reporter interface {
	HandleMessage(message string)
	Run()
}

type BaseReporter struct {
	cfg config.Config
}

func (r *BaseReporter) GetScoreForText(text string) (success bool, score float64, maxScore float64, penalties []scoring.ScoringPenalty, info *data.GithubRepoInfo) {
	match := GITHUB_URL_MATCHER.MatchString(text)
	if match {
		gh := helpers.GithubHelper{
			Config: r.cfg.Github,
		}
		ghInfo, _ := gh.GetRepositoryData(text)

		info = &ghInfo
		score, maxScore, penalties = scoring.GetTotalScore(info, r.cfg.Scorers)
		success = true
		return
	} else {
		return false, -1, -1, nil, nil
	}
}
