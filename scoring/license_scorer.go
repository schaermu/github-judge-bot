package scoring

import (
	"fmt"

	"github.com/mitchellh/go-spdx"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/helpers"
)

type LicenseScorer struct {
	data   helpers.GithubRepoInfo
	config config.LicenseConfig
}

func (s LicenseScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
	// the license of the project is either checked against a whitelist or against all osi approved licenses from spdx
	scoreChange := 0.0
	if len(s.config.ValidLicenseIds) == 0 {
		// TODO: cache this list
		list, _ := spdx.List()
		for _, spdxLic := range list.Licenses {
			if spdxLic.OSIApproved && !spdxLic.Deprecated {
				s.config.ValidLicenseIds = append(s.config.ValidLicenseIds, spdxLic.ID)
			}
		}
	}

	scoreChange = s.config.MaxPenalty
	for _, id := range s.config.ValidLicenseIds {
		if id == s.data.LicenseId {
			scoreChange = s.config.MaxPenalty
			break
		}
	}

	if scoreChange > 0 {
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			Reason: fmt.Sprintf("No valid license found: %s", s.data.License),
			Amount: scoreChange,
		})
	}

	return currentScore, penalties
}
