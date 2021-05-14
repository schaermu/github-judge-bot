package helpers

import (
	"context"
	"errors"
	"regexp"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
)

type GithubRepoInfo struct {
	OrgName        string
	RepositoryName string
	Stars          int
}

type GithubHelper struct {
	Config config.GithubConfig
}

func (gh GithubHelper) GetRepositoryData(repoUrl string) (info GithubRepoInfo, err error) {
	tp := github.BasicAuthTransport{
		Username: gh.Config.Username,
		Password: gh.Config.PersonalAccessToken,
	}
	client := github.NewClient(tp.Client())
	ctx := context.Background()

	org, repo, err := ExtractInfoFromUrl(repoUrl)
	if err != nil {
		return
	}

	repoInfo, _, err := client.Repositories.Get(ctx, org, repo)
	if err != nil {
		return
	}

	return GithubRepoInfo{
		OrgName:        org,
		RepositoryName: repo,
		Stars:          *repoInfo.StargazersCount,
	}, nil
}

func ExtractInfoFromUrl(repoUrl string) (org string, repo string, err error) {
	r, _ := regexp.Compile("github.com/([^/]+)/([^/<>]+)")
	matches := r.FindStringSubmatch(repoUrl)
	if len(matches) < 2 {
		return "", "", errors.New("could not determine organization and/or repository name from url")
	}
	return matches[1], matches[2], nil
}
