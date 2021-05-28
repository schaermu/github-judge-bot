package helpers

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"sync"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

// GithubURLRegex provides a regular expression to detect and split github repository url's.
const GithubURLRegex string = "github.com/([^/]+)/([^/<>]+)"

// GithubHelper provides helper functions to fetch data from Github's API.
type GithubHelper struct {
	Config config.GithubConfig
}

// GetRepositoryData fetches the latest data from Github API for a specific repository url.
func (gh GithubHelper) GetRepositoryData(repoURL string) (info data.GithubRepoInfo, err error) {
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

	org, repo, err := ExtractInfoFromURL(repoURL)
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
		LicenseID:      repoInfo.GetLicense().GetSPDXID(),
	}

	var waitgroup sync.WaitGroup
	waitgroup.Add(3)

	// fetch additional data
	go getRepoIssues(&info, client, &waitgroup)
	go getCommitActivity(&info, client, &waitgroup)
	go getContributors(&info, client, &waitgroup)

	waitgroup.Wait()

	return
}

func getContributors(info *data.GithubRepoInfo, client *github.Client, waitgroup *sync.WaitGroup) {
	contributors, _, _ := client.Repositories.ListContributors(context.Background(), info.OrgName, info.RepositoryName, &github.ListContributorsOptions{Anon: "1"})
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

// ExtractInfoFromURL gets the organisation/user and the repository slug from a Github URL.
func ExtractInfoFromURL(repoURL string) (org string, repo string, err error) {
	matches := regexp.MustCompile(GithubURLRegex).FindStringSubmatch(repoURL)
	if len(matches) < 2 {
		return "", "", errors.New("could not determine organization and/or repository name from url")
	}
	return matches[1], matches[2], nil
}
