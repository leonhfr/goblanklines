package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// checkTestHelperRule reports a missing blank line after a run of
// t.Parallel()/t.Helper() calls, unless the run ends its block.
func checkTestHelperRule(pass *analysis.Pass, stmts []ast.Stmt, tb *types.Interface) {
	checkAfterGroupRule(pass, stmts, func(stmt ast.Stmt) bool {
		return isTestHelperCall(pass, stmt, tb)
	}, "missing newline after test helper calls")
}

func isTestHelperCall(pass *analysis.Pass, stmt ast.Stmt, tb *types.Interface) bool {
	sel, ok := calledSelector(stmt)
	if !ok {
		return false
	}

	if sel.Sel.Name != "Parallel" && sel.Sel.Name != "Helper" {
		return false
	}

	return implementsTB(pass, sel.X, tb)
}
