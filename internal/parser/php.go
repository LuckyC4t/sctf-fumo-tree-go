package parser

import (
	"errors"
	"github.com/VKCOM/php-parser/pkg/ast"
	"github.com/VKCOM/php-parser/pkg/conf"
	phperr "github.com/VKCOM/php-parser/pkg/errors"
	"github.com/VKCOM/php-parser/pkg/parser"
	"github.com/VKCOM/php-parser/pkg/version"
	"os"
)

func ParsePHP(file string) (ast.Vertex, []error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, []error{err}
	}

	var parserErrors []*phperr.Error
	errorHandler := func(e *phperr.Error) {
		parserErrors = append(parserErrors, e)
	}

	rootNode, err := parser.Parse(content, conf.Config{
		Version: &version.Version{
			Major: 8,
			Minor: 0,
		},
		ErrorHandlerFunc: errorHandler,
	})

	if len(parserErrors) != 0 {
		errs := make([]error, len(parserErrors))
		for i, e := range parserErrors {
			errs[i] = errors.New(e.String())
		}
		return nil, errs
	}

	return rootNode, nil
}
