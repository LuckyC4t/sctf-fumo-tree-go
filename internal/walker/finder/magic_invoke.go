package finder

import (
	"encoding/base64"
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/visitor"
	"strings"
)

type MagicInvokeFinder struct {
	visitor.Null

	b64Key     string
	FindInvoke bool
}

func NewMagicInvokeFinder(b64key string) *MagicInvokeFinder {
	return &MagicInvokeFinder{
		b64Key: base64.StdEncoding.EncodeToString([]byte(b64key)),
	}
}

func (f *MagicInvokeFinder) ExprFunctionCall(n *ast.ExprFunctionCall) {
	funcName := string(n.Function.(*ast.Name).Parts[0].(*ast.NamePart).Value)
	if funcName == "base64_decode" {
		firstArg := n.Args[0].(*ast.Argument).Expr
		b64Str := string(firstArg.(*ast.ScalarString).Value)
		b64Str = strings.ReplaceAll(b64Str, "'", "")

		if b64Str == f.b64Key {
			f.FindInvoke = true
			return
		}
	}
}
