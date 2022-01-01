package data

import (
	"fmt"
	"github.com/VKCOM/php-parser/pkg/ast"
)

var ClassMethodIndex = make(map[string][]PHPClassMethod)

type PHPClassMethod struct {
	Name      string
	ClassName string
	Method    *ast.StmtClassMethod
}

func (pcm PHPClassMethod) String() string {
	return fmt.Sprintf("%s|%s", pcm.ClassName, pcm.Name)
}
