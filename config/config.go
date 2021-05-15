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
}

type StarsScoringConfig struct {
	MinStars   int     `mapstructure:"min_stars"`
	MaxPenalty float64 `mapstructure:"max_penalty"`
}

type IssuesScoringConfig struct {
	OpenClosedRatio float64 `mapstructure:"open_closed_ratio"`
	MaxPenalty      float64 `mapstructure:"max_penalty"`
}

type CommitActivityConfig struct {
	WeeklyInactivityPenalty float64 `mapstructure:"weekly_penalty"`
	MaxPenalty              float64 `mapstructure:"max_penalty"`
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
