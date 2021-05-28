package data

import "github.com/google/go-github/v35/github"

// GithubRepoInfo contains information from Github's API about a repository.
type GithubRepoInfo struct {
	OrgName        string
	RepositoryName string
	Stars          int
	License        string
	LicenseID      string
	Issues         []*github.Issue
	CommitActivity []*github.WeeklyCommitActivity
	Contributors   []*github.Contributor
}
