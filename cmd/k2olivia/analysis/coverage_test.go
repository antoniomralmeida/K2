package analysis

import (
	"testing"

	"github.com/antoniomralmeida/k2/cmd/k2olivia/modules"
	"github.com/antoniomralmeida/k2/cmd/k2olivia/util"
)

func TestGetModuleCoverage(t *testing.T) {
	defaultModules = modules.GetModules("en")

	coverage := getModuleCoverage("en")

	if len(coverage.NotCovered) != 0 || coverage.Coverage != 100 {
		t.Errorf("GetModuleCoverage() failed.")
	}
}

func TestGetIntentCoverage(t *testing.T) {
	defaultIntents = GetIntents("en")

	coverage := getIntentCoverage("en")

	if len(coverage.NotCovered) != 0 || coverage.Coverage != 100 {
		t.Errorf("GetIntentCoverage() failed.")
	}
}

func TestGetMessageCoverage(t *testing.T) {
	defaultMessages = util.GetMessages("en")

	coverage := getIntentCoverage("en")

	if len(coverage.NotCovered) != 0 || coverage.Coverage != 100 {
		t.Errorf("GetIntentCoverage() failed.")
	}
}
