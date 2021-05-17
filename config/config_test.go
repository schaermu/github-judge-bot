package config

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	appConf, err := New(bytes.NewBuffer(GetTestConfig()))

	if err != nil {
		t.Errorf("Failed to unmarshal config: %s", err)
	}

	if len(appConf.Scorers) != 5 {
		t.Errorf("Wrong number of scorers was read: %d", len(appConf.Scorers))
	}
}
