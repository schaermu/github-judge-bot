package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractInfoFromUrl(t *testing.T) {
	org, repo, err := ExtractInfoFromURL("https://github.com/foo/bar")
	assert.Nil(t, err)
	assert.Equal(t, "foo", org)
	assert.Equal(t, "bar", repo)
}

func TestExtractInfoFromUrlNonGithub(t *testing.T) {
	org, repo, err := ExtractInfoFromURL("https://foobar.org/test/me")
	assert.NotNil(t, err)
	assert.Empty(t, org)
	assert.Empty(t, repo)
}

func TestExtractInfoFromUrlInvalid(t *testing.T) {
	org, repo, err := ExtractInfoFromURL("THIS_IS_SPARTA")
	assert.NotNil(t, err)
	assert.Empty(t, org)
	assert.Empty(t, repo)
}

func TestGetRepositoryData(t *testing.T) {
	ghHelper := GithubHelper{}
	info, err := ghHelper.GetRepositoryData("https://github.com/schaermu/github-judge-bot")

	assert.Nil(t, err)
	assert.Equal(t, "github-judge-bot", info.RepositoryName)
	assert.GreaterOrEqual(t, 1, len(info.Contributors))
}
