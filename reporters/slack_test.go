package reporters

import (
	"fmt"
	"testing"

	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/schaermu/go-github-judge-bot/scoring"
	"github.com/stretchr/testify/assert"
)

func TestGetSlackMessageColors(t *testing.T) {
	color, icon := getSlackMessageColorAndIcon(10, 10)
	assert.Equal(t, "good", color)
	assert.Equal(t, ":+1::skin-tone-2:", icon)
}

func TestGetSlackMessageColorsWarn(t *testing.T) {
	color, icon := getSlackMessageColorAndIcon(10, 7)
	assert.Equal(t, "warning", color)
	assert.Equal(t, ":warning:", icon)
}
func TestGetSlackMessageColorsDanger(t *testing.T) {
	color, icon := getSlackMessageColorAndIcon(10, 2)
	assert.Equal(t, "danger", color)
	assert.Equal(t, ":exclamation:", icon)
}

func TestBuildSlackResponse(t *testing.T) {
	penalties := []scoring.Penalty{}
	msgBlocks := buildSlackResponse(&data.GithubRepoInfo{OrgName: "foo", RepositoryName: "bar"}, 10, 10, penalties)
	assert.Len(t, msgBlocks, 2)
}

func TestBuildSlackResponsePenalties(t *testing.T) {
	penalties := []scoring.Penalty{{Reason: "Test reason", Amount: 3}}
	msgBlocks := buildSlackResponse(&data.GithubRepoInfo{OrgName: "foo", RepositoryName: "bar"}, 7, 10, penalties)
	assert.Len(t, msgBlocks, 3)
}

func TestBuildSlackError(t *testing.T) {
	msgBlocks := buildSlackError(&data.GithubRepoInfo{OrgName: "foo", RepositoryName: "bar"}, fmt.Errorf("foo bar failed"))
	assert.Len(t, msgBlocks, 3)
}
