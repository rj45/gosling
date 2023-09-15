package main

import (
	"os"

	"github.com/rj45/gosling/arch/aarch64"
	"github.com/rj45/gosling/codegen"
	"github.com/rj45/gosling/parser"
)

func main() {
	src := []byte(os.Args[1])
	parser := parser.New(src)
	ast := parser.Parse()

	gen := codegen.New(ast, &aarch64.Assembly{})
	gen.Generate()
}
