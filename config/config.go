package config

import (
	"io"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var defaultConfig = []byte(`
debug: false
scorers:
  - name: stars
    max_penalty: 2.0
    settings:
        min_stars: 800
  - name: issues
    max_penalty: 2.0
    settings:
        closed_open_ratio: 0.2 # maximum of open tickets per closed ones
  - name: commit-activity
    max_penalty: 3.0
    settings:
        weekly_penalty: 0.1
  - name: contributors
    max_penalty: 1.0
    settings:
        min_contributors: 3
  - name: license
    max_penalty: 2.0
`)

func GetDefaultConfig() []byte {
	return defaultConfig
}

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

type Config struct {
	Slack     SlackConfig    `mapstructure:"slack"`
	Github    GithubConfig   `mapstructure:"github"`
	Scorers   []ScorerConfig `mapstructure:"scorers"`
	DebugMode bool           `mapstructure:"debug"`
}

type ScorerConfig struct {
	Name       string            `mapstructure:"name"`
	MaxPenalty float64           `mapstructure:"max_penalty"`
	Enabled    bool              `mapstructure:"enabled"`
	Settings   map[string]string `mapstructure:"settings"`
}

func (s ScorerConfig) Get(key string) string {
	if val, found := s.Settings[key]; found {
		return val
	}
	return ""
}

func (s ScorerConfig) GetFloat64(key string) float64 {
	val := s.Get(key)
	if val != "" {
		return cast.ToFloat64(val)
	}
	return 0
}

func (s ScorerConfig) GetInt(key string) int {
	val := s.Get(key)
	if val != "" {
		return cast.ToInt(val)
	}
	return 0
}

func (s ScorerConfig) GetSlice(key string) []string {
	val := s.Get(key)
	if val != "" {
		return cast.ToStringSlice(val)
	}
	return []string{}
}

func New(rdr io.Reader) (config Config, err error) {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.ReadConfig(rdr)
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
