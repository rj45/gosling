package compile

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/codegen"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/semantics"
)

func Compile(file *ast.File, asm codegen.Assembly) []error {
	parser := parser.New(file)
	ast, errs := parser.Parse()
	if errs != nil {
		return errs
	}

	symtab, errs := semantics.NewTypeChecker(ast).Check(ast.Root())
	if errs != nil {
		return errs
	}

	gen := codegen.New(ast, symtab, asm)
	gen.Generate()

	return nil
}
