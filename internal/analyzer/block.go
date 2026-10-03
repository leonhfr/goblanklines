package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// checkBlockRule checks a statement list for missing blank lines after
// block statements (if, for, switch, select, etc.).
func checkBlockRule(pass *analysis.Pass, astFile *ast.File, stmts []ast.Stmt) {
	for i := 0; i < len(stmts)-1; i++ {
		checkBlockStatementPair(pass, astFile, stmts[i], stmts[i+1])
	}

	if len(stmts) > 0 {
		checkBlockLastStatement(pass, astFile, stmts[len(stmts)-1])
	}
}

// checkSwitchLikeClauses checks that case/comm clauses in switch, type
// switch, and select statements are separated by blank lines.
func checkSwitchLikeClauses(pass *analysis.Pass, astFile *ast.File, stmts []ast.Stmt) {
	var clauses []ast.Stmt

	for _, stmt := range stmts {
		switch stmt.(type) {
		case *ast.CaseClause, *ast.CommClause:
			clauses = append(clauses, stmt)
		}
	}

	for i := 0; i < len(clauses)-1; i++ {
		checkClauseGap(pass, astFile, clauses[i], clauses[i+1].Pos())
	}
}

func clauseBody(stmt ast.Stmt) []ast.Stmt {
	switch c := stmt.(type) {
	case *ast.CaseClause:
		return c.Body

	case *ast.CommClause:
		return c.Body
	}

	return nil
}

func checkClauseGap(pass *analysis.Pass, astFile *ast.File, current ast.Stmt, nextPos token.Pos) {
	body := clauseBody(current)
	if len(body) == 0 {
		return
	}

	lastStmtEnd := body[len(body)-1].End()

	file := pass.Fset.File(lastStmtEnd)
	if file == nil {
		return
	}

	lastStmtLine := file.Line(lastStmtEnd)

	foundComment := checkCommentBetween(pass, astFile, file, lastStmtEnd, lastStmtLine, nextPos, "missing newline after case block")
	if !foundComment && file.Line(nextPos) == lastStmtLine+1 {
		pass.Report(afterFix(pass, lastStmtEnd, "missing newline after case block"))
	}
}

// checkBlockStatementPair checks if there's proper spacing between two
// consecutive statements.
func checkBlockStatementPair(pass *analysis.Pass, astFile *ast.File, current, next ast.Stmt) {
	if isDeferStmt(current) && isDeferStmt(next) {
		return
	}

	if isErrorCheckIfStmt(pass, current) && isDeferStmt(next) {
		return
	}

	if !needsNewlineAfter(current) {
		return
	}

	blockEnd := getBlockEnd(current)
	if blockEnd == token.NoPos {
		return
	}

	file := pass.Fset.File(blockEnd)
	if file == nil {
		return
	}

	blockEndLine := file.Line(blockEnd)

	foundComment := checkCommentBetween(pass, astFile, file, blockEnd, blockEndLine, next.Pos(), "missing newline after block statement")
	if !foundComment && file.Line(next.Pos()) == blockEndLine+1 {
		pass.Report(afterFix(pass, blockEnd, "missing newline after block statement"))
	}
}

// checkBlockLastStatement checks if the last statement in a list has proper
// spacing before any trailing comment.
func checkBlockLastStatement(pass *analysis.Pass, astFile *ast.File, lastStmt ast.Stmt) {
	if !needsNewlineAfter(lastStmt) {
		return
	}

	blockEnd := getBlockEnd(lastStmt)
	if blockEnd == token.NoPos {
		return
	}

	file := pass.Fset.File(blockEnd)
	if file == nil {
		return
	}

	blockEndLine := file.Line(blockEnd)

	for _, cg := range astFile.Comments {
		if cg.Pos() <= blockEnd {
			continue
		}

		commentLine := file.Line(cg.Pos())
		if commentLine == blockEndLine {
			continue
		}

		if commentLine == blockEndLine+1 {
			pass.Report(afterFix(pass, blockEnd, "missing newline after block statement"))
		}

		break
	}
}

// checkCommentBetween reports message if a comment right after afterLine
// precedes beforePos, and reports whether such a comment exists at all.
func checkCommentBetween(pass *analysis.Pass, astFile *ast.File, file *token.File, afterPos token.Pos, afterLine int, beforePos token.Pos, message string) bool {
	for _, cg := range astFile.Comments {
		pos := cg.Pos()
		if pos <= afterPos || pos >= beforePos {
			continue
		}

		line := file.Line(pos)
		if line == afterLine {
			continue
		}

		if line == afterLine+1 {
			pass.Report(afterFix(pass, afterPos, message))
		}

		return true
	}

	return false
}

// needsNewlineAfter determines if a statement needs a newline after it.
func needsNewlineAfter(stmt ast.Stmt) bool {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		return true

	case *ast.ForStmt, *ast.RangeStmt:
		return true

	case *ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
		return true

	case *ast.AssignStmt:
		return extractAssignFuncLit(s) != nil

	case *ast.DeclStmt:
		return extractDeclFuncLit(s) != nil

	case *ast.DeferStmt:
		return true
	}

	return false
}

// getBlockEnd returns the end position of a block statement's body.
func getBlockEnd(stmt ast.Stmt) token.Pos {
	if pos := blockBodyEnd(stmt); pos != token.NoPos {
		return pos
	}

	return funcLitBodyEnd(stmt)
}

// blockBodyEnd handles the block-statement cases of getBlockEnd.
func blockBodyEnd(stmt ast.Stmt) token.Pos {
	if pos, ok := ifStmtEnd(stmt); ok {
		return pos
	}

	switch s := stmt.(type) {
	case *ast.BlockStmt:
		return s.End()

	case *ast.DeferStmt:
		return s.End()
	}

	if body := stmtBody(stmt); body != nil {
		return body.End()
	}

	return token.NoPos
}

// ifStmtEnd handles the if-statement case of getBlockEnd. ok is false when
// stmt is not an *ast.IfStmt.
func ifStmtEnd(stmt ast.Stmt) (token.Pos, bool) {
	s, ok := stmt.(*ast.IfStmt)
	if !ok {
		return token.NoPos, false
	}

	if s.Else != nil {
		return getBlockEnd(s.Else), true
	}

	if s.Body != nil {
		return s.Body.End(), true
	}

	return token.NoPos, true
}

// stmtBody returns the body of the statement types that carry a single
// *ast.BlockStmt body, or nil otherwise.
func stmtBody(stmt ast.Stmt) *ast.BlockStmt {
	switch s := stmt.(type) {
	case *ast.ForStmt:
		return s.Body

	case *ast.RangeStmt:
		return s.Body

	case *ast.SwitchStmt:
		return s.Body

	case *ast.TypeSwitchStmt:
		return s.Body

	case *ast.SelectStmt:
		return s.Body
	}

	return nil
}

// funcLitBodyEnd handles the function-literal cases of getBlockEnd.
func funcLitBodyEnd(stmt ast.Stmt) token.Pos {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if funcLit := extractAssignFuncLit(s); funcLit != nil && funcLit.Body != nil {
			return funcLit.Body.End()
		}

	case *ast.DeclStmt:
		if funcLit := extractDeclFuncLit(s); funcLit != nil && funcLit.Body != nil {
			return funcLit.Body.End()
		}
	}

	return token.NoPos
}

// extractAssignFuncLit extracts a non-invoked function literal assigned by s.
func extractAssignFuncLit(s *ast.AssignStmt) *ast.FuncLit {
	for _, expr := range s.Rhs {
		if funcLit := asFuncLit(expr); funcLit != nil {
			return funcLit
		}
	}

	return nil
}

// extractDeclFuncLit extracts a non-invoked function literal declared by s.
func extractDeclFuncLit(s *ast.DeclStmt) *ast.FuncLit {
	genDecl, ok := s.Decl.(*ast.GenDecl)
	if !ok {
		return nil
	}

	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		for _, value := range valueSpec.Values {
			if funcLit := asFuncLit(value); funcLit != nil {
				return funcLit
			}
		}
	}

	return nil
}

// asFuncLit returns expr if it is a function literal that is not
// immediately invoked (an invoked literal ends with ')', not '}').
func asFuncLit(expr ast.Expr) *ast.FuncLit {
	if funcLit, ok := expr.(*ast.FuncLit); ok {
		return funcLit
	}

	return nil
}

// isErrorCheckIfStmt checks if an if statement matches "if <error> != nil".
func isErrorCheckIfStmt(pass *analysis.Pass, stmt ast.Stmt) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok {
		return false
	}

	binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok || binaryExpr.Op != token.NEQ {
		return false
	}

	return isErrNotNilPattern(pass, binaryExpr.X, binaryExpr.Y) || isErrNotNilPattern(pass, binaryExpr.Y, binaryExpr.X)
}

// isErrNotNilPattern checks if x is a variable implementing the error
// interface and y is the nil identifier.
func isErrNotNilPattern(pass *analysis.Pass, x, y ast.Expr) bool {
	ident, ok := x.(*ast.Ident)
	if !ok {
		return false
	}

	nilIdent, ok := y.(*ast.Ident)
	if !ok || nilIdent.Name != "nil" {
		return false
	}

	if pass.TypesInfo == nil {
		return false
	}

	typ := pass.TypesInfo.TypeOf(ident)
	if typ == nil {
		return false
	}

	return implementsError(typ)
}

// implementsError checks if a type implements the error interface.
func implementsError(typ types.Type) bool {
	errorObj := types.Universe.Lookup("error")
	if errorObj == nil {
		return false
	}

	errorType, ok := errorObj.Type().Underlying().(*types.Interface)
	if !ok {
		return false
	}

	if types.Implements(typ, errorType) {
		return true
	}

	if _, ok := typ.(*types.Pointer); !ok {
		return types.Implements(types.NewPointer(typ), errorType)
	}

	return false
}

// isDeferStmt checks if a statement is a defer statement.
func isDeferStmt(stmt ast.Stmt) bool {
	_, ok := stmt.(*ast.DeferStmt)
	return ok
}
