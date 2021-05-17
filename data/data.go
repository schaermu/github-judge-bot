package data

import "github.com/google/go-github/v35/github"

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
