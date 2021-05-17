package scoring

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
	"github.com/stretchr/testify/assert"
)

func getGithubTestData() helpers.GithubRepoInfo {
	return helpers.GithubRepoInfo{
		License:        "Not existing dummy license",
		LicenseId:      "NOT_EXISTING",
		Stars:          1,
		OrgName:        "foo",
		RepositoryName: "bar",
		Issues:         getTestIssueData(0.2, 10, 3),
		CommitActivity: getTestCommitActivityData(2),
		Contributors:   getTestContributorData(3),
	}
}

func getTestConfig() (config.Config, error) {
	appConf, err := config.New(bytes.NewBuffer(config.GetTestConfig()))
	return appConf, err
}

func TestCreateScorer(t *testing.T) {
	testData := getGithubTestData()
	config, err := getTestConfig()
	if err != nil {
		panic(err)
	}

	for _, config := range config.Scorers {
		scorer := CreateScorer(testData, config)
		assert.Equal(t, strcase.ToCamel(config.Name+"Scorer"), reflect.TypeOf(scorer).Name())
	}
}

func TestGetTotalScore(t *testing.T) {
	testData := getGithubTestData()
	config, err := getTestConfig()
	if err != nil {
		panic(err)
	}

	_, maxScore := CreateScorerMap(testData, config.Scorers)

	score, penalties := GetTotalScore(testData, config.Scorers)
	assert.NotEqual(t, maxScore, score)
	assert.NotEmpty(t, penalties)
}
