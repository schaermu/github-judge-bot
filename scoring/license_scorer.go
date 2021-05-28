package scoring

import (
	"fmt"

	"github.com/mitchellh/go-spdx"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

type LicenseScorer struct {
	data   *data.GithubRepoInfo
	config config.ScorerConfig
}

func (s LicenseScorer) GetScore(currentScore float64, penalties []ScoringPenalty) (float64, []ScoringPenalty) {
	// the license of the project is either checked against a whitelist or against all osi approved licenses from spdx
	scoreChange := 0.0
	validIds := s.config.GetSlice("valid_license_ids")
	if len(validIds) == 0 {
		// TODO: cache this list
		list, _ := spdx.List()
		for _, spdxLic := range list.Licenses {
			if spdxLic.OSIApproved && !spdxLic.Deprecated {
				validIds = append(validIds, spdxLic.ID)
			}
		}
	}

	scoreChange = s.config.MaxPenalty
	for _, id := range validIds {
		if id == s.data.LicenseID {
			scoreChange = 0
			break
		}
	}

	if scoreChange > 0 {
		currentScore -= scoreChange

		penalties = append(penalties, ScoringPenalty{
			ScorerName: "License",
			Reason:     fmt.Sprintf("No valid license found: %s", s.data.License),
			Amount:     scoreChange,
		})
	}

	return currentScore, penalties
}
