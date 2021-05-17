package helpers

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"sync"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
)

type GithubRepoInfo struct {
	OrgName        string
	RepositoryName string
	Stars          int
	License        string
	LicenseId      string
	Issues         []*github.Issue
	CommitActivity []*github.WeeklyCommitActivity
	Contributors   []*github.ContributorStats
}

type GithubHelper struct {
	Config config.GithubConfig
}

func (gh GithubHelper) GetRepositoryData(repoUrl string) (info GithubRepoInfo, err error) {
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

	info = GithubRepoInfo{
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

func getContributorStats(info *GithubRepoInfo, client *github.Client, waitgroup *sync.WaitGroup) {
	contributors, _, _ := client.Repositories.ListContributorsStats(context.Background(), info.OrgName, info.RepositoryName)
	info.Contributors = contributors
	waitgroup.Done()
}

func getCommitActivity(info *GithubRepoInfo, client *github.Client, waitgroup *sync.WaitGroup) {
	commitActivity, _, _ := client.Repositories.ListCommitActivity(context.Background(), info.OrgName, info.RepositoryName)
	info.CommitActivity = commitActivity
	waitgroup.Done()
}

func getRepoIssues(info *GithubRepoInfo, client *github.Client, waitgroup *sync.WaitGroup) {
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
