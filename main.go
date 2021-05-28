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

	args := os.Args[1:]

	if len(args) > 0 {
		// assume stdin as source for url, use stdout reporter
		stdoutReporter := reporters.NewStdoutReporter(&cfg)
		stdoutReporter.HandleMessage(os.Args[1:][0])
	} else if cfg.Slack.AppToken != "" {
		// start handling events coming in from slack
		slackReporter := reporters.NewSlackReporter(&cfg)
		go slackReporter.HandleMessage("")
		slackReporter.Run()
	} else {
		// print usage
		println("Usage: github-judge [URL]")
		println("")
		println("Note: to start github-judge in bot-mode, make sure to configure it properly.")
	}
}
