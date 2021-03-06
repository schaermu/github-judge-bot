package scoring

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/iancoleman/strcase"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/stretchr/testify/assert"
)

func getGithubTestData() *data.GithubRepoInfo {
	return &data.GithubRepoInfo{
		License:        "Not existing dummy license",
		LicenseID:      "NOT_EXISTING",
		Stars:          1,
		OrgName:        "foo",
		RepositoryName: "bar",
		Issues:         getTestIssueData(0.2, 10, 3),
		CommitActivity: getTestCommitActivityData(2),
		Contributors:   getTestContributorData(3),
	}
}

func getTestConfig() (config.Config, error) {
	appConf, err := config.New(bytes.NewBuffer(config.GetDefaultConfig()))
	return appConf, err
}

func TestCreateScorer(t *testing.T) {
	testData := getGithubTestData()
	config, err := getTestConfig()
	if err != nil {
		panic(err)
	}

	for _, sc := range config.Scorers {
		scorer := CreateScorer(testData, sc)
		assert.Equal(t, strcase.ToCamel(sc.Name+"Scorer"), reflect.TypeOf(scorer).Name())
	}
}

func TestGetTotalScore(t *testing.T) {
	testData := getGithubTestData()
	config, err := getTestConfig()
	if err != nil {
		panic(err)
	}

	summary := GetTotalScore(testData, config.Scorers)
	assert.NotEqual(t, summary.MaxScore, summary.Score)
	assert.NotEmpty(t, summary.Penalties)
}
