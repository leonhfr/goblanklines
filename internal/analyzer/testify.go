package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

const (
	testifyAssertPkg  = "github.com/stretchr/testify/assert"
	testifyRequirePkg = "github.com/stretchr/testify/require"
)

// checkTestifyRule reports a missing blank line after a run of testify
// assert/require calls, unless the run ends its block.
func checkTestifyRule(pass *analysis.Pass, stmts []ast.Stmt) {
	checkAfterGroupRule(pass, stmts, func(stmt ast.Stmt) bool {
		return isTestifyCall(pass, stmt)
	}, "missing newline after testify assertion group")
}

// isTestifyCall reports whether stmt is a bare call to a function or
// *Assertions method defined in the testify assert or require package.
func isTestifyCall(pass *analysis.Pass, stmt ast.Stmt) bool {
	sel, ok := calledSelector(stmt)
	if !ok {
		return false
	}

	fn, ok := pass.TypesInfo.ObjectOf(sel.Sel).(*types.Func)
	if !ok {
		return false
	}

	pkg := fn.Pkg()
	if pkg == nil {
		return false
	}

	return pkg.Path() == testifyAssertPkg || pkg.Path() == testifyRequirePkg
}
