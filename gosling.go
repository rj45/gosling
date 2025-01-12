package main

import (
	"fmt"
	"os"

	"github.com/rj45/gosling/arch/aarch64"
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/compile"
)

func main() {
	src := []byte(os.Args[1])
	file := ast.NewFile("test.gos", src)
	asm := &aarch64.Assembly{Out: os.Stdout}
	errs := compile.Compile(file, asm)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}
