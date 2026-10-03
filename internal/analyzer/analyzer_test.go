package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestBlockRule(t *testing.T) {
	t.Parallel()

	runRule(t, ruleBlock, "blockstatements", "structliterals", "comments", "caseclauses", "deferpattern", "anonymousfuncs")
}

func TestBlockRuleWithFixes(t *testing.T) {
	t.Parallel()

	runRuleWithFixes(t, ruleBlock, "blockstatements", "structliterals", "comments", "caseclauses", "deferpattern", "anonymousfuncs")
}

func TestTestHelperRule(t *testing.T) {
	t.Parallel()

	runRule(t, ruleTestHelper, "testhelper")
}

func TestTestHelperRuleWithFixes(t *testing.T) {
	t.Parallel()

	runRuleWithFixes(t, ruleTestHelper, "testhelper")
}

func TestTestifyRule(t *testing.T) {
	t.Parallel()

	runRule(t, ruleTestify, "testify")
}

func TestTestifyRuleWithFixes(t *testing.T) {
	t.Parallel()

	runRuleWithFixes(t, ruleTestify, "testify")
}

func TestTestRunRule(t *testing.T) {
	t.Parallel()

	runRule(t, ruleTestRun, "testrun")
}

func TestTestRunRuleWithFixes(t *testing.T) {
	t.Parallel()

	runRuleWithFixes(t, ruleTestRun, "testrun")
}

func TestFuncDeclRule(t *testing.T) {
	t.Parallel()

	runRule(t, ruleFuncDecl, "funcdecl")
}

func TestFuncDeclRuleWithFixes(t *testing.T) {
	t.Parallel()

	runRuleWithFixes(t, ruleFuncDecl, "funcdecl")
}

// TestCombined exercises all rules together against a single fixture that
// mixes the patterns each rule targets.
func TestCombined(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, New(Settings{}), "combined")
}

func TestCombinedWithFixes(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, New(Settings{}), "combined")
}

// TestRulesSetting verifies that the "rules" setting actually restricts
// which rules run: a fixture with both a testhelper and a testify violation
// should only report the testhelper violation when rules=[testhelper].
func TestRulesSetting(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, New(Settings{Rules: []string{"testhelper"}}), "rulessubset")
}

// runRule runs only the given rule, so other rules' diagnostics don't
// affect each rule's fixtures.
func runRule(t *testing.T, rule string, pkgs ...string) {
	t.Helper()

	testdata := analysistest.TestData()

	for _, pkg := range pkgs {
		t.Run(pkg, func(t *testing.T) {
			t.Parallel()

			analysistest.Run(t, testdata, New(Settings{Rules: []string{rule}}), pkg)
		})
	}
}

func runRuleWithFixes(t *testing.T, rule string, pkgs ...string) {
	t.Helper()

	testdata := analysistest.TestData()

	for _, pkg := range pkgs {
		t.Run(pkg, func(t *testing.T) {
			t.Parallel()

			analysistest.RunWithSuggestedFixes(t, testdata, New(Settings{Rules: []string{rule}}), pkg)
		})
	}
}
