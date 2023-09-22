package main

import (
	"fmt"
	"os"

	"github.com/rj45/gosling/arch/aarch64"
	"github.com/rj45/gosling/codegen"
	"github.com/rj45/gosling/parser"
)

func main() {
	src := []byte(os.Args[1])
	parser := parser.New(src)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	gen := codegen.New(ast, &aarch64.Assembly{Out: os.Stdout})
	gen.Generate()
}
