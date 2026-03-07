package analyzer

import (
	"go/ast"
	"go/types"
)

func isZapLogger(sel *ast.SelectorExpr, info *types.Info) bool {
	obj := info.Uses[sel.Sel]
	if obj == nil {
		return false
	}
	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	path := pkg.Path()
	return path == "go.uber.org/zap"
}

func isSlogCall(sel *ast.SelectorExpr, info *types.Info) bool {
	obj := info.Uses[sel.Sel]
	if obj == nil {
		return false
	}

	pkg := obj.Pkg()
	if pkg == nil {
		return false
	}

	path := pkg.Path()
	return path == "log/slog"
}

func isLogMethod(name string) bool {
	switch name {
	case "Info", "Error", "Warn", "Debug", "Panic", "Fatal",
		"DPanic", "Sync":
		return true
	}
	return false
}
