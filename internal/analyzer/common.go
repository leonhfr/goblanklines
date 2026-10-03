package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// blankLineBetween reports whether at least one blank line separates the
// line containing end from the line containing start.
func blankLineBetween(pass *analysis.Pass, end, start token.Pos) bool {
	file := pass.Fset.File(end)
	if file == nil {
		return true
	}

	return file.Line(start) > file.Line(end)+1
}

// afterFix builds a diagnostic that inserts a blank line right after the
// line containing pos.
func afterFix(pass *analysis.Pass, pos token.Pos, message string) analysis.Diagnostic {
	file := pass.Fset.File(pos)
	if file == nil {
		return analysis.Diagnostic{Pos: pos, Message: message}
	}

	line := file.Line(pos)

	var insertPos token.Pos
	if line < file.LineCount() {
		insertPos = file.LineStart(line + 1)
	} else {
		insertPos = token.Pos(file.Base() + file.Size())
	}

	return analysis.Diagnostic{
		Pos:     pos,
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message:   "Insert blank line",
			TextEdits: []analysis.TextEdit{{Pos: insertPos, End: insertPos, NewText: []byte("\n")}},
		}},
	}
}

// beforeFix builds a diagnostic that inserts a blank line right before the
// line containing pos.
func beforeFix(pass *analysis.Pass, pos token.Pos, message string) analysis.Diagnostic {
	file := pass.Fset.File(pos)
	if file == nil {
		return analysis.Diagnostic{Pos: pos, Message: message}
	}

	insertPos := file.LineStart(file.Line(pos))

	return analysis.Diagnostic{
		Pos:     pos,
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message:   "Insert blank line",
			TextEdits: []analysis.TextEdit{{Pos: insertPos, End: insertPos, NewText: []byte("\n")}},
		}},
	}
}

// calledSelector extracts the selector of a bare statement of the form
// x.Sel(...).
func calledSelector(stmt ast.Stmt) (*ast.SelectorExpr, bool) {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return nil, false
	}

	call, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	return sel, true
}

// lookupTestingTB finds the testing.TB interface type among the analyzed
// package's direct imports, or nil if "testing" is not imported.
func lookupTestingTB(pass *analysis.Pass) *types.Interface {
	for _, imp := range pass.Pkg.Imports() {
		if imp.Path() != "testing" {
			continue
		}

		obj := imp.Scope().Lookup("TB")
		if obj == nil {
			return nil
		}

		iface, ok := obj.Type().Underlying().(*types.Interface)
		if !ok {
			return nil
		}

		return iface
	}

	return nil
}

// implementsTB reports whether expr's type implements testing.TB.
func implementsTB(pass *analysis.Pass, expr ast.Expr, tb *types.Interface) bool {
	if tb == nil {
		return false
	}

	typ := pass.TypesInfo.TypeOf(expr)

	return typ != nil && types.Implements(typ, tb)
}

// checkAfterGroupRule reports a missing blank line after each maximal run of
// consecutive statements matching isMember, unless the run ends its block.
func checkAfterGroupRule(pass *analysis.Pass, stmts []ast.Stmt, isMember func(ast.Stmt) bool, message string) {
	for i := 0; i < len(stmts); {
		if !isMember(stmts[i]) {
			i++
			continue
		}

		for i < len(stmts) && isMember(stmts[i]) {
			i++
		}

		last := i - 1
		if last == len(stmts)-1 {
			continue
		}

		if !blankLineBetween(pass, stmts[last].End(), stmts[last+1].Pos()) {
			pass.Report(afterFix(pass, stmts[last].End(), message))
		}
	}
}

// checkIsolatedGroupRule reports missing blank lines around each maximal
// run of isMember statements, unless the run starts or ends its block.
func checkIsolatedGroupRule(pass *analysis.Pass, stmts []ast.Stmt, isMember func(ast.Stmt) bool, beforeMsg, afterMsg string) {
	for i := 0; i < len(stmts); {
		if !isMember(stmts[i]) {
			i++
			continue
		}

		start := i

		for i < len(stmts) && isMember(stmts[i]) {
			i++
		}

		last := i - 1

		if start > 0 && !blankLineBetween(pass, stmts[start-1].End(), stmts[start].Pos()) {
			pass.Report(beforeFix(pass, stmts[start].Pos(), beforeMsg))
		}

		if last < len(stmts)-1 && !blankLineBetween(pass, stmts[last].End(), stmts[last+1].Pos()) {
			pass.Report(afterFix(pass, stmts[last].End(), afterMsg))
		}
	}
}
