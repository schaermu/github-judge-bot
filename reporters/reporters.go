package reporters

import (
	"fmt"
	"regexp"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/schaermu/go-github-judge-bot/helpers"
	"github.com/schaermu/go-github-judge-bot/scoring"
)

// Reporter provides the interface all reporters must follow.
type Reporter interface {
	HandleMessage(message string)
	Run()
}

// BaseReporter provides a base functionality for reporters.
type BaseReporter struct {
	cfg config.Config
}

func (r *BaseReporter) getScoreForText(text string) (success bool, summary scoring.ScoringSummary, info *data.GithubRepoInfo, err error) {
	match := regexp.MustCompile(helpers.GITHUB_URL_REGEX).MatchString(text)
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
	}
	return false, summary, nil, fmt.Errorf("%s does not contain a valid github.com url", text)
}
