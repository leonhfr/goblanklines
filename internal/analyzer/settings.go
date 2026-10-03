package analyzer

// Settings configures which rules the analyzer runs.
type Settings struct {
	// Rules restricts the analyzer to the listed rule keys. An empty list
	// (the default) enables every rule.
	Rules []string `json:"rules"`
}

const (
	ruleBlock      = "block"
	ruleTestHelper = "testhelper"
	ruleTestify    = "testify"
	ruleTestRun    = "testrun"
	ruleFuncDecl   = "funcdecl"
)

var allRules = []string{ruleBlock, ruleTestHelper, ruleTestify, ruleTestRun, ruleFuncDecl}

// enabledRules resolves the settings into the effective set of active rules.
func (s Settings) enabledRules() map[string]bool {
	keys := s.Rules
	if len(keys) == 0 {
		keys = allRules
	}

	enabled := make(map[string]bool, len(keys))

	for _, r := range keys {
		enabled[r] = true
	}

	return enabled
}
