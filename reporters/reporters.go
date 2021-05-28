package reporters

import (
	"fmt"
	"regexp"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/schaermu/go-github-judge-bot/helpers"
	"github.com/schaermu/go-github-judge-bot/scoring"
)

var GITHUB_URL_MATCHER = regexp.MustCompile(helpers.GITHUB_URL_REGEX)

type Reporter interface {
	HandleMessage(message string)
	Run()
}

type BaseReporter struct {
	cfg config.Config
}

func (r *BaseReporter) getScoreForText(text string) (success bool, summary scoring.ScoringSummary, info *data.GithubRepoInfo, err error) {
	match := GITHUB_URL_MATCHER.MatchString(text)
	if match {
		gh := helpers.GithubHelper{
			Config: r.cfg.Github,
		}
		ghInfo, githubErr := gh.GetRepositoryData(text)
		if githubErr != nil {
			err = githubErr
			success = false
			return
		}

		info = &ghInfo
		summary = scoring.GetTotalScore(info, r.cfg.Scorers)
		success = true
		return
	} else {
		return false, summary, nil, fmt.Errorf("%s does not contain a valid github.com url", text)
	}
}
