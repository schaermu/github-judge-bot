package config

import (
	"github.com/spf13/viper"
)

type SlackConfig struct {
	BotToken   string `mapstructure:"bot_token"`
	AppToken   string `mapstructure:"app_token"`
	SigningKey string `mapstructure:"signing_key"`
	Debug      bool   `mapstructure:"debug"`
}

type GithubConfig struct {
	Username            string `mapstructure:"username"`
	PersonalAccessToken string `mapstructure:"access_token"`
}

type ScoringConfig struct {
	MaxScore       float64              `mapstructure:"max_score"`
	Stars          StarsScoringConfig   `mapstructure:"stars"`
	Issues         IssuesScoringConfig  `mapstructure:"issues"`
	CommitActivity CommitActivityConfig `mapstructure:"activity"`
	Contributors   ContributorsConfig   `mapstructure:"contributors"`
	License        LicenseConfig        `mapstructure:"license"`
}

type StarsScoringConfig struct {
	MaxPenalty float64 `mapstructure:"max_penalty"`
	Enabled    bool    `mapstructure:"enabled"`
	MinStars   int     `mapstructure:"min_stars"`
}

type IssuesScoringConfig struct {
	MaxPenalty      float64 `mapstructure:"max_penalty"`
	Enabled         bool    `mapstructure:"enabled"`
	ClosedOpenRatio float64 `mapstructure:"closed_open_ratio"`
}

type CommitActivityConfig struct {
	MaxPenalty              float64 `mapstructure:"max_penalty"`
	Enabled                 bool    `mapstructure:"enabled"`
	WeeklyInactivityPenalty float64 `mapstructure:"weekly_penalty"`
}

type ContributorsConfig struct {
	MaxPenalty      float64 `mapstructure:"max_penalty"`
	Enabled         bool    `mapstructure:"enabled"`
	MinContributors int     `mapstructure:"min_contributors"`
}

type LicenseConfig struct {
	MaxPenalty      float64  `mapstructure:"max_penalty"`
	Enabled         bool     `mapstructure:"enabled"`
	ValidLicenseIds []string `mapstructure:"valid_license_ids"`
}

type Config struct {
	Slack     SlackConfig   `mapstructure:"slack"`
	Github    GithubConfig  `mapstructure:"github"`
	Score     ScoringConfig `mapstructure:"scoring"`
	DebugMode bool          `mapstructure:"debug"`
}

func New() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
