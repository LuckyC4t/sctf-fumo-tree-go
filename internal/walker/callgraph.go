package walker

import (
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/data"
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/util"
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/walker/finder"
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
	"github.com/VKCOM/php-parser/pkg/visitor/traverser"
	"strings"
)

type CallGraphWalker struct {
	visitor.Null

	VariableStatus map[string]bool
	CallOut        []data.PHPClassMethod
}

func NewCallGraphWalker() *CallGraphWalker {
	return &CallGraphWalker{
		VariableStatus: make(map[string]bool),
		CallOut:        []data.PHPClassMethod{},
	}
}

func (w *CallGraphWalker) StmtClassMethod(n *ast.StmtClassMethod) {
	// 获取参数，默认参数变量可控
	for _, p := range n.Params {
		varName := util.GetVariableName(p.(*ast.Parameter).Var.(*ast.ExprVariable).Name.(*ast.Identifier).Value)
		w.VariableStatus[varName] = true
	}
}

func (w *CallGraphWalker) ExprAssign(n *ast.ExprAssign) {
	// 简单的赋值: $b = foo($a) | $b = $a
	right := util.GetVariableName(n.Var.(*ast.ExprVariable).Name.(*ast.Identifier).Value)
	switch left := n.Expr.(type) {
	case *ast.ExprFunctionCall:
		funcName := string(left.Function.(*ast.Name).Parts[0].(*ast.NamePart).Value)
		if data.Sanitizer[funcName] {
			w.VariableStatus[right] = false
		}
	case *ast.ExprVariable:
		varName := util.GetVariableName(left.Name.(*ast.Identifier).Value)
		w.VariableStatus[right] = w.VariableStatus[varName]
	}
}

// 处理 $this->aaa->bbb(ccc)
func (w *CallGraphWalker) ExprMethodCall(n *ast.ExprMethodCall) {
	for _, arg := range n.Args {
		if a, ok := arg.(*ast.Argument); ok {
			if e, ok := a.Expr.(*ast.ExprVariable); ok {
				varName := string(e.Name.(*ast.Identifier).Value)
				varName = strings.ReplaceAll(varName, "$", "") // 去除 $ 标记
				if polluted, has := w.VariableStatus[varName]; !has || (has && !polluted) {
					// 变量不受控制
					return
				}
			}
		}
	}

	targetName := string(n.Method.(*ast.Identifier).Value)
	targetMethods := w.GetTargetMethods(targetName)
	w.CallOut = append(w.CallOut, targetMethods...)
}

func (w *CallGraphWalker) ExprFunctionCall(n *ast.ExprFunctionCall) {
	funcName := string(n.Function.(*ast.Name).Parts[0].(*ast.NamePart).Value)
	if funcName == "extract" {
		// 获取变量对应的字符串的值，做为后续 call_user_func 的目标
		firstArg := n.Args[0].(*ast.Argument).Expr.(*ast.ExprArray)
		targetName := string(firstArg.Items[0].(*ast.ExprArrayItem).Val.(*ast.ScalarString).Value)
		targetName = strings.ReplaceAll(targetName, "'", "")

		targetMethods := w.GetTargetMethods(targetName)
		w.CallOut = append(w.CallOut, targetMethods...)
	} else if funcName == "call_user_func" {
		// @call_user_func($this->WHB5xkK7, ['LUlnpp' => $RwGAFc8G])
		firstArg := n.Args[0].(*ast.Argument).Expr
		// 与另一种 call_user_func 区分开
		if _, ok := firstArg.(*ast.ExprPropertyFetch); ok {
			secondArg := n.Args[1].(*ast.Argument).Expr.(*ast.ExprArray)
			b64Key := string(secondArg.Items[0].(*ast.ExprArrayItem).Key.(*ast.ScalarString).Value)
			b64Key = strings.ReplaceAll(b64Key, "'", "")

			varName := util.GetVariableName(secondArg.Items[0].(*ast.ExprArrayItem).
				Val.(*ast.ExprVariable).Name.(*ast.Identifier).Value)

			if w.VariableStatus[varName] {
				// 寻找 invoke
				invokeFinder := finder.NewMagicInvokeFinder(b64Key)
				for _, invoke := range data.ClassMethodIndex["__invoke"] {
					traverser.NewTraverser(invokeFinder).Traverse(invoke.Method)
					if invokeFinder.FindInvoke {
						w.CallOut = append(w.CallOut, invoke)

						invokeFinder.FindInvoke = false
					}
				}
			}
		}
	}
}

// self util

func (w *CallGraphWalker) GetTargetMethods(targetName string) []data.PHPClassMethod {
	var res []data.PHPClassMethod
	// 从缓存中查照对应名称的方法
	if methods, has := data.ClassMethodIndex[targetName]; has {
		res = append(res, methods...)
	} else { // 寻找调用 call_user_func 中 对应变量名与当前所寻找的方法名一致的 __call
		callFinder := finder.NewMagicCallFinder(targetName)
		for _, call := range data.ClassMethodIndex["__call"] {
			traverser.NewTraverser(callFinder).Traverse(call.Method)
			if callFinder.FindTarget {
				res = append(res, call)

				callFinder.FindTarget = false
			}
		}
	}

	return res
}
