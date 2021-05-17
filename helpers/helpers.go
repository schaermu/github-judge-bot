package helpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/schaermu/go-github-judge-bot/scoring"
	"github.com/slack-go/slack"
)

type GithubHelper struct {
	Config config.GithubConfig
}

func (gh GithubHelper) GetRepositoryData(repoUrl string) (info data.GithubRepoInfo, err error) {
	tp := http.DefaultClient
	if gh.Config.Username != "" && gh.Config.PersonalAccessToken != "" {
		basicAuth := github.BasicAuthTransport{
			Username: gh.Config.Username,
			Password: gh.Config.PersonalAccessToken,
		}
		tp = basicAuth.Client()
	}

	client := github.NewClient(tp)
	ctx := context.Background()

	org, repo, err := ExtractInfoFromUrl(repoUrl)
	if err != nil {
		return
	}

	repoInfo, _, err := client.Repositories.Get(ctx, org, repo)
	if err != nil {
		return
	}

	info = data.GithubRepoInfo{
		OrgName:        org,
		RepositoryName: repo,
		Stars:          repoInfo.GetStargazersCount(),
		License:        repoInfo.GetLicense().GetName(),
		LicenseId:      repoInfo.GetLicense().GetSPDXID(),
	}

	var waitgroup sync.WaitGroup
	waitgroup.Add(3)

	// fetch additional data
	go getRepoIssues(&info, client, &waitgroup)
	go getCommitActivity(&info, client, &waitgroup)
	go getContributorStats(&info, client, &waitgroup)

	waitgroup.Wait()

	return
}

func getContributorStats(info *data.GithubRepoInfo, client *github.Client, waitgroup *sync.WaitGroup) {
	contributors, _, _ := client.Repositories.ListContributorsStats(context.Background(), info.OrgName, info.RepositoryName)
	info.Contributors = contributors
	waitgroup.Done()
}

func getCommitActivity(info *data.GithubRepoInfo, client *github.Client, waitgroup *sync.WaitGroup) {
	commitActivity, _, _ := client.Repositories.ListCommitActivity(context.Background(), info.OrgName, info.RepositoryName)
	info.CommitActivity = commitActivity
	waitgroup.Done()
}

func getRepoIssues(info *data.GithubRepoInfo, client *github.Client, waitgroup *sync.WaitGroup) {
	issues, _, _ := client.Issues.ListByRepo(context.Background(), info.OrgName, info.RepositoryName, &github.IssueListByRepoOptions{
		Sort:  "updated",
		State: "all",
	})
	info.Issues = issues
	waitgroup.Done()
}

func ExtractInfoFromUrl(repoUrl string) (org string, repo string, err error) {
	r, _ := regexp.Compile("github.com/([^/]+)/([^/<>]+)")
	matches := r.FindStringSubmatch(repoUrl)
	if len(matches) < 2 {
		return "", "", errors.New("could not determine organization and/or repository name from url")
	}
	return matches[1], matches[2], nil
}

func GetSlackMessageColorAndIcon(score float64, maxScore float64) (color string, icon string) {
	if maxScore/100*score < .4 {
		return "danger", ":exclamation:"
	}
	if maxScore/100*score < .8 {
		return "warning", ":warning:"
	}
	return "good", ":+1::skin-tone-2:"
}

func BuildSlackResponse(org string, repository string, score float64, maxScore float64, penalties []scoring.ScoringPenalty) []slack.MsgOption {
	messageColor, messageIcon := GetSlackMessageColorAndIcon(score, maxScore)

	// build default message
	res := []slack.MsgOption{
		slack.MsgOptionIconEmoji(messageIcon),
		slack.MsgOptionText(fmt.Sprintf("Analysis of `%s/%s` complete, final score is *%.2f/10.00* points!", org, repository, score), false),
	}

	// append penalty attachment containing details
	if len(penalties) > 0 {
		penaltyOutput := ""
		for _, penalty := range penalties {
			penaltyOutput += fmt.Sprintf("-*%.2f* _%s_\n", penalty.Amount, penalty.Reason)
		}

		attachment := slack.MsgOptionAttachments(
			slack.Attachment{
				Color:      messageColor,
				MarkdownIn: []string{"text"},
				Text:       penaltyOutput,
				Pretext:    "The following reasons lead to penalties",
			},
		)

		res = append(res, attachment)
	}

	return res
}
