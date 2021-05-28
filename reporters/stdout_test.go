package reporters

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func getTestConfig() (config.Config, error) {
	appConf, err := config.New(bytes.NewBuffer(config.GetDefaultConfig()))
	return appConf, err
}

func TestStdoutReporter(t *testing.T) {
	cfg, _ := getTestConfig()

	stdoutReporter := NewStdoutReporter(&cfg)
	output := captureOutput(func() {
		stdoutReporter.HandleMessage("github.com/schaermu/github-judge-bot")
	})

	assert.Contains(t, "github-judge-bot", output)
}
