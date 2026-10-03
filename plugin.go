// Package goblanklines registers the goblanklines linter as a
// golangci-lint module plugin.
//
// Build a custom golangci-lint binary that includes it by listing this
// module in .custom-gcl.yml and running `golangci-lint custom`; then enable
// the linter under linters.settings.custom in .golangci.yml.
package goblanklines

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/leonhfr/goblanklines/internal/analyzer"
)

func init() {
	register.Plugin("goblanklines", newPlugin)
}

type plugin struct {
	settings Settings
}

func newPlugin(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	return &plugin{settings: s}, nil
}

// BuildAnalyzers returns the single analyzer this plugin provides.
func (p *plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.New(p.settings)}, nil
}

// GetLoadMode requires full type information: the testhelper, testify, and
// testrun rules all type-check call receivers and resolved functions.
func (p *plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
