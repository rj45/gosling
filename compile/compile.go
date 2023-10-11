package compile

import (
	"github.com/rj45/gosling/codegen"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/semantics"
	"github.com/rj45/gosling/token"
)

func Compile(file *token.File, asm codegen.Assembly) []error {
	parser := parser.New(file)
	ast, errs := parser.Parse()
	if errs != nil {
		return errs
	}

	tc := semantics.NewTypeChecker(ast)

	symtab, errs := tc.Check(ast.Root())
	if errs != nil {
		return errs
	}

	gen := codegen.New(ast, symtab, tc.Universe(), asm)
	gen.Generate()

	return nil
}
