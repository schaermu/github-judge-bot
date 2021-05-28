package reporters

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/schaermu/go-github-judge-bot/config"
)

type StdoutReporter struct {
	Reporter
	BaseReporter
}

func NewStdoutReporter(cfg *config.Config) StdoutReporter {
	return StdoutReporter{
		BaseReporter: BaseReporter{cfg: *cfg},
	}
}

func (sr *StdoutReporter) Run() {
	// noop
}

func (sr *StdoutReporter) HandleMessage(message string) {
	if isScored, summary, info, err := sr.getScoreForText(message); isScored && err == nil {
		// pretty print result to stdout
		println()
		println(fmt.Sprintf("Judgement of of %s/%s complete, it scored %.2f/%.2f points.", info.OrgName, info.RepositoryName, summary.Score, summary.MaxScore))
		println()

		if len(summary.Penalties) > 0 {
			t := table.NewWriter()
			t.SetStyle(table.StyleRounded)
			t.SetColumnConfigs([]table.ColumnConfig{
				{
					Name:        "Penalty",
					Align:       text.AlignRight,
					AlignHeader: text.AlignRight,
					AlignFooter: text.AlignRight,
				},
			})
			t.SetOutputMirror(os.Stdout)
			t.AppendHeader(table.Row{"Scorer", "Reason", "Penalty"})

			for _, penalty := range summary.Penalties {
				t.AppendRow(table.Row{penalty.ScorerName, penalty.Reason, fmt.Sprintf("-%.2f", penalty.Amount)})
				t.AppendSeparator()
			}
			t.AppendFooter(table.Row{"", "TOTAL", fmt.Sprintf("-%.2f", summary.TotalPenalties)})
			t.Render()
		}
	} else if err != nil {
		println(fmt.Sprintf("Judgement of '%s' failed: %s", message, err.Error()))
	}
}
