package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// checkFuncDeclRule reports a missing blank line between two consecutive
// func/method declarations, unless both are one-liners.
func checkFuncDeclRule(pass *analysis.Pass, file *ast.File) {
	var prev *ast.FuncDecl

	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			prev = nil
			continue
		}

		if prev != nil && (!isOneLinerFunc(pass, prev) || !isOneLinerFunc(pass, fd)) {
			pos := funcDeclStart(fd)
			if !blankLineBetween(pass, prev.End(), pos) {
				pass.Report(beforeFix(pass, pos, "missing newline between function declarations"))
			}
		}

		prev = fd
	}
}

func funcDeclStart(fd *ast.FuncDecl) token.Pos {
	if fd.Doc != nil {
		return fd.Doc.Pos()
	}

	return fd.Pos()
}

func isOneLinerFunc(pass *analysis.Pass, fd *ast.FuncDecl) bool {
	file := pass.Fset.File(fd.Pos())
	if file == nil {
		return false
	}

	return file.Line(fd.Pos()) == file.Line(fd.End())
}
