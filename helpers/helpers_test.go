package helpers

import (
	"testing"

	"github.com/schaermu/go-github-judge-bot/scoring"
	"github.com/stretchr/testify/assert"
)

func TestExtractInfoFromUrl(t *testing.T) {
	org, repo, err := ExtractInfoFromUrl("https://github.com/foo/bar")
	assert.Nil(t, err)
	assert.Equal(t, "foo", org)
	assert.Equal(t, "bar", repo)
}

func TestExtractInfoFromUrlNonGithub(t *testing.T) {
	org, repo, err := ExtractInfoFromUrl("https://foobar.org/test/me")
	assert.NotNil(t, err)
	assert.Empty(t, org)
	assert.Empty(t, repo)
}

func TestExtractInfoFromUrlInvalid(t *testing.T) {
	org, repo, err := ExtractInfoFromUrl("THIS_IS_SPARTA")
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

func TestGetSlackMessageColors(t *testing.T) {
	color, icon := GetSlackMessageColorAndIcon(10, 10)
	assert.Equal(t, "good", color)
	assert.Equal(t, ":+1::skin-tone-2:", icon)
}

func TestGetSlackMessageColorsWarn(t *testing.T) {
	color, icon := GetSlackMessageColorAndIcon(10, 7)
	assert.Equal(t, "warning", color)
	assert.Equal(t, ":warning:", icon)
}
func TestGetSlackMessageColorsDanger(t *testing.T) {
	color, icon := GetSlackMessageColorAndIcon(10, 2)
	assert.Equal(t, "danger", color)
	assert.Equal(t, ":exclamation:", icon)
}

func TestBuildSlackResponse(t *testing.T) {
	penalties := []scoring.ScoringPenalty{}

	msgBlocks := BuildSlackResponse("foo", "bar", 10, 10, penalties)

	assert.Len(t, msgBlocks, 2)
}

func TestBuildSlackResponsePenalties(t *testing.T) {
	penalties := []scoring.ScoringPenalty{{Reason: "Test reason", Amount: 3}}

	msgBlocks := BuildSlackResponse("foo", "bar", 7, 10, penalties)

	assert.Len(t, msgBlocks, 3)
}
