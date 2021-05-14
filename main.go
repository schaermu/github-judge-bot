package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
	"github.com/schaermu/go-github-judge-bot/scoring"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println("Failed to load config, aborting...")
		return
	}

	api := slack.New(
		cfg.Slack.BotToken,
		slack.OptionDebug(cfg.Slack.Debug),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(cfg.Slack.AppToken),
	)

	client := socketmode.New(
		api,
		socketmode.OptionDebug(cfg.Slack.Debug),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	messageMatcher, _ := regexp.Compile("github.com/([^/]+)/([^/<>]+)")

	go func() {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("Connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Println("Connection failed. Retrying later...")
			case socketmode.EventTypeConnected:
				fmt.Println("Connected to Slack with Socket Mode.")
			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					fmt.Printf("Ignored %+v\n", evt)
					continue
				}

				fmt.Printf("Event received: %+v\n", eventsAPIEvent)

				client.Ack(*evt.Request)

				switch eventsAPIEvent.Type {
				case slackevents.CallbackEvent:
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						match := messageMatcher.Match([]byte(ev.Text))
						if match {
							gh := helpers.GithubHelper{
								Config: cfg.Github,
							}
							info, err := gh.GetRepositoryData(ev.Text)

							if err != nil {
								fmt.Printf("Error while fetching github info: %e", err)
							}

							score, penalties := scoring.GetTotalScore(info, cfg.Score)
							fmt.Printf("Final score: %f\n", score)
							for _, penalty := range penalties {
								fmt.Printf("%f: %s\n", penalty.Amount, penalty.Reason)
							}
						}
					}
				default:
					client.Debugf("unsupported Events API event received")
				}
			default:
				fmt.Fprintf(os.Stderr, "Unexpected event type received: %s\n", evt.Type)
			}
		}
	}()

	client.Run()
}
