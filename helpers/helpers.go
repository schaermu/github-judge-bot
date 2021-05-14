package helpers

import (
	"context"
	"errors"
	"regexp"

	"github.com/google/go-github/v35/github"
	"github.com/schaermu/go-github-judge-bot/config"
)

type GithubRepoInfo struct {
	Stars int
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

	org, user, err := ExtractInfoFromUrl(repoUrl)
	if err != nil {
		return
	}

	repo, _, err := client.Repositories.Get(ctx, org, user)
	if err != nil {
		return
	}

	return GithubRepoInfo{
		Stars: *repo.StargazersCount,
	}, nil
}

func ExtractInfoFromUrl(repoUrl string) (org string, user string, err error) {
	r, _ := regexp.Compile("github.com/([^/]+)/([^/<>]+)")
	matches := r.FindStringSubmatch(repoUrl)
	if len(matches) < 2 {
		return "", "", errors.New("could not determine organization and/or username from repository url")
	}
	return matches[1], matches[2], nil
}
