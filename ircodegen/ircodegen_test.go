package ircodegen

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/ir"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/semantics"
	"github.com/rj45/gosling/vm"
)

var tests = []struct {
	name   string
	input  string
	output int
}{
	{
		name:   "return zero",
		input:  `{return 0}`,
		output: 0,
	},
	{
		name:   "return constant",
		input:  `{return 42}`,
		output: 42,
	},
	{
		name:   "return simple expression",
		input:  `{return 1+2-2}`,
		output: 1,
	},
	{
		name:   "return addition",
		input:  `{return 5 + 20 - 4}`,
		output: 21,
	},
	{
		name:   "return addition with spaces",
		input:  `{return  12 + 34 - 5 }`,
		output: 41,
	},
	{
		name:   "return multiplication",
		input:  `{return 5+6*7}`,
		output: 47,
	},
	{
		name:   "return parentheses",
		input:  `{return 5*(9-6)}`,
		output: 15,
	},
	{
		name:   "return division",
		input:  `{return (3+5)/2}`,
		output: 4,
	},
	{
		name:   "return negative number",
		input:  `{return -10+20}`,
		output: 10,
	},
	{
		name:   "return double negative",
		input:  `{return - -10}`,
		output: 10,
	},
	{
		name:   "return triple negative",
		input:  `{return - - +10}`,
		output: 10,
	},
	{
		name:   "if statement false",
		input:  `{if false {return 1} else {return 0}}`,
		output: 0,
	},
	// TODO: fixme
	// {
	// 	name:   "if statement true",
	// 	input:  `{if true {return 1} else {return 0}}`,
	// 	output: 1,
	// },
}

func TestCodegenWithVirtualMachine(t *testing.T) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			input := "func main() int " + tt.input
			if strings.Contains(tt.input, "main()") {
				input = tt.input
			}
			file := ast.NewFile("test.gos", []byte(input))

			parser := parser.New(file)
			ast, errs := parser.Parse()
			for _, err := range errs {
				t.Errorf("Expected no error, but got\n%s", err)
			}

			symtab, errs := semantics.NewTypeChecker(ast).Check(ast.Root())
			for _, err := range errs {
				t.Errorf("Expected no error, but got\n%s", err)
			}

			translator := ir.NewSoNTranslator(ast)
			program := translator.Translate()
			asm := vm.NewAsm()

			codegen := New(program, symtab, asm)
			codegen.Generate()

			if len(errs) > 0 {
				for _, err := range errs {
					t.Errorf("Expected no error, but got\n%s", err)
				}
			}

			vm := vm.NewCPU(asm.Program)
			vm.Trace = true
			actual := vm.Run()
			if actual != tt.output {
				t.Errorf("Expected: %d; but got: %d", tt.output, actual)
			}
		})
	}
}
