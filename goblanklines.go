// Package goblanklines provides a linter that inserts blank lines
// according to a set of rules. See the internal/analyzer package for the
// rule implementations.
package goblanklines

import (
	"golang.org/x/tools/go/analysis"

	"github.com/leonhfr/goblanklines/internal/analyzer"
)

// Settings configures which rules the analyzer runs.
type Settings = analyzer.Settings

// New creates a new goblanklines analyzer with every rule enabled.
func New() *analysis.Analyzer {
	return analyzer.New(Settings{})
}
