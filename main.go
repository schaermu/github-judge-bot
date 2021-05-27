package main

import (
	"log"
	"os"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/reporters"
)

func main() {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	cfg, err := config.New(f)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		return
	}

	if cfg.Slack.AppToken != "" {
		// start handling events coming in from slack
		slackReporter := reporters.NewSlackReporter(&cfg)
		go slackReporter.HandleMessage("")
		slackReporter.Run()
	} else {
		// assume stdin as source for url, use terminal reporter

	}
}
