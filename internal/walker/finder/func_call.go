package finder

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
)

type FuncCallFinder struct {
	visitor.Null

	funcName     string
	FindFuncCall bool
}

func NewFuncCallFinder(funcName string) *FuncCallFinder {
	return &FuncCallFinder{funcName: funcName}
}

func (f *FuncCallFinder) ExprFunctionCall(n *ast.ExprFunctionCall) {
	funcName := string(n.Function.(*ast.Name).Parts[0].(*ast.NamePart).Value)
	if funcName == f.funcName {
		f.FindFuncCall = true
		return
	}
}
