package scoring

import (
	"fmt"

	"github.com/mitchellh/go-spdx"
	"github.com/schaermu/go-github-judge-bot/config"
	"github.com/schaermu/go-github-judge-bot/data"
)

// LicenseScorer provides a scoring based on the repository's license.
type LicenseScorer struct {
	data   *data.GithubRepoInfo
	config config.ScorerConfig
}

// GetScore calculates a score based on the licensing.
// If the repository does not provide a valid, OSI approved and non-deprecated license listed on SPDX, a penalty is applied.
// Optionally, the user can configure a custom whitelist in "valid_license_ids" for all license he want's to let through (this overrules the SPDX index).
func (s LicenseScorer) GetScore(currentScore float64, penalties []Penalty) (float64, []Penalty) {
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

		penalties = append(penalties, Penalty{
			ScorerName: "License",
			Reason:     fmt.Sprintf("No valid license found: %s", s.data.License),
			Amount:     scoreChange,
		})
	}

	return currentScore, penalties
}
