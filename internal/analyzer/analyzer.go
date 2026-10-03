// Package analyzer implements the goblanklines rules as a single
// golang.org/x/tools/go/analysis.Analyzer.
package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const doc = `check for missing blank lines

This linter enforces blank lines in several situations, each independently
toggled through the "rules" setting (default: all enabled):

- block: a blank line after block statements (if/for/switch/select/defer),
  ported from github.com/leonhfr/newline-after-block.
- testhelper: a blank line after a run of t.Parallel()/t.Helper() calls.
- testify: blank lines before and after a run of testify assert/require
  calls, isolating them from surrounding code.
- testrun: blank lines before and after a run of t.Run(...) calls.
- funcdecl: a blank line between consecutive function/method declarations,
  unless both are one-liners.

The analyzer provides automatic fix suggestions that insert the required
blank lines.`

type goblanklines struct {
	enabled map[string]bool
}

// New creates a new goblanklines analyzer configured by s.
func New(s Settings) *analysis.Analyzer {
	g := &goblanklines{enabled: s.enabledRules()}

	return &analysis.Analyzer{
		Name: "goblanklines",
		Doc:  doc,
		Run:  g.run,
	}
}

func (g *goblanklines) run(pass *analysis.Pass) (any, error) {
	tb := lookupTestingTB(pass)

	for _, file := range pass.Files {
		if g.enabled[ruleFuncDecl] {
			checkFuncDeclRule(pass, file)
		}

		ast.Inspect(file, func(n ast.Node) bool {
			g.inspect(pass, file, tb, n)
			return true
		})
	}

	//nolint:nilnil // analysis.Analyzer.Run's contract: nil error means success with no shared result.
	return nil, nil
}

func (g *goblanklines) inspect(pass *analysis.Pass, file *ast.File, tb *types.Interface, n ast.Node) {
	switch x := n.(type) {
	case *ast.BlockStmt:
		g.checkStmtList(pass, file, tb, x.List)

	case *ast.CaseClause:
		g.checkStmtList(pass, file, tb, x.Body)

	case *ast.CommClause:
		g.checkStmtList(pass, file, tb, x.Body)

	default:
		if g.enabled[ruleBlock] {
			if body := switchLikeBody(n); body != nil {
				checkSwitchLikeClauses(pass, file, body)
			}
		}
	}
}

// switchLikeBody returns the clause list of a switch, type switch, or
// select statement, or nil if n is none of those or has no body.
func switchLikeBody(n ast.Node) []ast.Stmt {
	switch x := n.(type) {
	case *ast.SwitchStmt:
		if x.Body != nil {
			return x.Body.List
		}

	case *ast.TypeSwitchStmt:
		if x.Body != nil {
			return x.Body.List
		}

	case *ast.SelectStmt:
		if x.Body != nil {
			return x.Body.List
		}
	}

	return nil
}

func (g *goblanklines) checkStmtList(pass *analysis.Pass, file *ast.File, tb *types.Interface, stmts []ast.Stmt) {
	if g.enabled[ruleBlock] {
		checkBlockRule(pass, file, stmts)
	}

	if g.enabled[ruleTestHelper] {
		checkTestHelperRule(pass, stmts, tb)
	}

	if g.enabled[ruleTestify] {
		checkTestifyRule(pass, stmts)
	}

	if g.enabled[ruleTestRun] {
		checkTestRunRule(pass, stmts, tb)
	}
}
