package reporters

import (
	"fmt"
	"log"
	"os"

	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
	"github.com/schaermu/go-github-judge-bot/scoring"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// SlackReporter sends judge results to the configured slack team.
type SlackReporter struct {
	Reporter
	BaseReporter
	client *socketmode.Client
	api    *slack.Client
}

// NewSlackReporter creates a new SlackReporter instance based on the config.
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

// Run starts listening for new messages on the Slack socketmode client.
func (sr *SlackReporter) Run() {
	sr.client.Run()
}

// HandleMessage will react to all messages that are pushed through the Slack socketmode client.
// NOTE: always pass an empty string as a message, this parameter is being ignored!
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
					if isScored, summary, info, err := sr.getScoreForText(ev.Text); isScored && err == nil {
						msgBlocks := buildSlackResponse(info, summary.Score, summary.MaxScore, summary.Penalties)
						sr.api.PostMessage(ev.Channel, msgBlocks...)
					} else if err != nil {
						sr.api.PostMessage(ev.Channel, buildSlackError(info, err)...)
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

func getSlackMessageColorAndIcon(score float64, maxScore float64) (color string, icon string) {
	if maxScore/100*score < .4 {
		return "danger", ":exclamation:"
	}
	if maxScore/100*score < .8 {
		return "warning", ":warning:"
	}
	return "good", ":+1::skin-tone-2:"
}

func buildSlackError(repoInfo *data.GithubRepoInfo, err error) []slack.MsgOption {
	var messageColor = "danger"
	var messageIcon = ":exclamation:"

	return []slack.MsgOption{
		slack.MsgOptionIconEmoji(messageIcon),
		slack.MsgOptionText(fmt.Sprintf("Analysis of `%s/%s` failed!", repoInfo.OrgName, repoInfo.RepositoryName), false),
		slack.MsgOptionAttachments(
			slack.Attachment{
				Color:      messageColor,
				MarkdownIn: []string{"text"},
				Text:       err.Error(),
			},
		),
	}
}

func buildSlackResponse(repoInfo *data.GithubRepoInfo, score float64, maxScore float64, penalties []scoring.ScoringPenalty) []slack.MsgOption {
	messageColor, messageIcon := getSlackMessageColorAndIcon(score, maxScore)

	// build default message
	res := []slack.MsgOption{
		slack.MsgOptionIconEmoji(messageIcon),
		slack.MsgOptionText(fmt.Sprintf("Judgement of `%s/%s` complete, it scored *%.2f/%.2f* points!", repoInfo.OrgName, repoInfo.RepositoryName, score, maxScore), false),
	}

	// append penalty attachment containing details
	if len(penalties) > 0 {
		penaltyOutput := ""
		for _, penalty := range penalties {
			penaltyOutput += fmt.Sprintf("*-%.2f* _%s_\n", penalty.Amount, penalty.Reason)
		}

		attachment := slack.MsgOptionAttachments(
			slack.Attachment{
				Color:      messageColor,
				MarkdownIn: []string{"text"},
				Text:       penaltyOutput,
				Pretext:    "The following reasons lead to penalties:",
			},
		)

		res = append(res, attachment)
	}

	return res
}
