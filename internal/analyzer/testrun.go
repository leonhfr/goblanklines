package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// checkTestRunRule reports missing blank lines around a run of t.Run(...)
// calls, unless the run starts or ends its block.
func checkTestRunRule(pass *analysis.Pass, stmts []ast.Stmt, tb *types.Interface) {
	checkIsolatedGroupRule(pass, stmts, func(stmt ast.Stmt) bool {
		return isTestRunCall(pass, stmt, tb)
	}, "missing newline before t.Run calls", "missing newline after t.Run calls")
}

func isTestRunCall(pass *analysis.Pass, stmt ast.Stmt, tb *types.Interface) bool {
	sel, ok := calledSelector(stmt)
	if !ok || sel.Sel.Name != "Run" {
		return false
	}

	return implementsTB(pass, sel.X, tb)
}
