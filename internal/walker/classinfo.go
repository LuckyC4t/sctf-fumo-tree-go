package walker

import (
	"github.com/LuckyC4t/sctf-fumo-tree-go/internal/data"
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
)

type ClassInfoWalker struct {
	visitor.Null

	PHPClasses []data.PHPClass
}

func NewClassInfoWalker() *ClassInfoWalker {
	return &ClassInfoWalker{}
}

func (w *ClassInfoWalker) StmtClass(n *ast.StmtClass) {
	className := string(n.Name.(*ast.Identifier).Value)

	// 提取类中的方法
	var phpClsMethods []data.PHPClassMethod
	for _, stmt := range n.Stmts {
		if classMethod, ok := stmt.(*ast.StmtClassMethod); ok {
			methodName := string(classMethod.Name.(*ast.Identifier).Value)
			phpClsMethod := data.PHPClassMethod{
				Name:      methodName,
				ClassName: className,
				Method:    classMethod,
			}
			phpClsMethods = append(phpClsMethods, phpClsMethod)
		}
	}

	w.PHPClasses = append(w.PHPClasses,
		data.PHPClass{
			Name:    className,
			Methods: phpClsMethods,
		})
}
