package finder

import (
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
	"strings"
)

type MagicCallFinder struct {
	visitor.Null

	TargetName string
	FindTarget bool
}

func NewMagicCallFinder(target string) *MagicCallFinder {
	return &MagicCallFinder{TargetName: target}
}

func (f *MagicCallFinder) ExprFunctionCall(n *ast.ExprFunctionCall) {
	funcName := string(n.Function.(*ast.Name).Parts[0].(*ast.NamePart).Value)
	if funcName == "call_user_func" {
		firstArg := n.Args[0].(*ast.Argument).Expr
		// call_user_func([$this->cgKM3ZI4, $ICaiygpY2z], ...$value)
		// 第一个参数为数组
		if arrayItem, ok := firstArg.(*ast.ExprArray); ok {
			// 获取数组中第二个 item 的变量名称
			varName := string(arrayItem.Items[1].(*ast.ExprArrayItem).
				Val.(*ast.ExprVariable).Name.(*ast.Identifier).Value)
			varName = strings.ReplaceAll(varName, "$", "")

			// 找到目标
			if varName == f.TargetName {
				f.FindTarget = true
				return
			}
		}
	}
}
