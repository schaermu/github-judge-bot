package reporters

import (
	"log"
	"os"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackReporter struct {
	Reporter
	BaseReporter
	client *socketmode.Client
	api    *slack.Client
}

func NewSlackReporter(cfg *config.Config) SlackReporter {
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

	return SlackReporter{
		BaseReporter: BaseReporter{cfg: *cfg},
		client:       client,
		api:          api,
	}
}

func (sr *SlackReporter) Run() {
	sr.client.Run()
}

func (sr *SlackReporter) HandleMessage(message string) {
	for evt := range sr.client.Events {
		switch evt.Type {
		case socketmode.EventTypeConnecting:
			log.Println("Connecting to Slack with Socket Mode...")
		case socketmode.EventTypeConnectionError:
			log.Println("Connection failed. Retrying later...")
		case socketmode.EventTypeConnected:
			log.Println("Connected to Slack with Socket Mode.")
		case socketmode.EventTypeHello:
			log.Println("Slack sent back hello, handshake complete.")
		case socketmode.EventTypeEventsAPI:
			eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				log.Printf("Ignored %+v\n", evt)
				continue
			}

			log.Printf("Event received: %+v\n", eventsAPIEvent)

			sr.client.Ack(*evt.Request)

			switch eventsAPIEvent.Type {
			case slackevents.CallbackEvent:
				innerEvent := eventsAPIEvent.InnerEvent
				switch ev := innerEvent.Data.(type) {
				case *slackevents.AppMentionEvent:
					if isScored, score, maxScore, penalties, info := sr.GetScoreForText(ev.Text); isScored {
						msgBlocks := helpers.BuildSlackResponse(info.OrgName, info.RepositoryName, score, maxScore, penalties)
						sr.api.PostMessage(ev.Channel, msgBlocks...)
					}
				}
			default:
				sr.client.Debugf("unsupported Events API event received")
			}
		default:
			log.Printf("Unexpected event type received: %s\n", evt.Type)
		}
	}
}
