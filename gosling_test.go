package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/rj45/gosling/arch/aarch64"
	"github.com/rj45/gosling/compile"
	"github.com/rj45/gosling/token"
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
		input:  `{if 0==1 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "if statement true",
		input:  `{if 42==42 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "not equal",
		input:  `{if 0!=1 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "not equal false",
		input:  `{if 42!=42 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "less than",
		input:  `{if 0<1 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "less than false",
		input:  `{if 1<1 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "less than false 2",
		input:  `{if 2<1 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "less than or equal",
		input:  `{if 0<=1 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "less than or equal true",
		input:  `{if 1<=1 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "less than or equal false",
		input:  `{if 2<=1 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "greater than",
		input:  `{if 1>0 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "greater than false",
		input:  `{if 1>1 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "greater than false 2",
		input:  `{if 1>2 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "greater than or equal",
		input:  `{if 1>=0 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "greater than or equal true",
		input:  `{if 1>=1 {return 1} else {return 0}}`,
		output: 1,
	},
	{
		name:   "greater than or equal false",
		input:  `{if 1>=2 {return 1} else {return 0}}`,
		output: 0,
	},
	{
		name:   "assign variable",
		input:  `{a:=3; return a}`,
		output: 3,
	},
	{
		name:   "add variables",
		input:  `{a:=3; z:=5; return a+z}`,
		output: 8,
	},
	{
		name:   "assign variable with letters",
		input:  `{foo:=3; return foo}`,
		output: 3,
	},
	{
		name:   "add variables with letters",
		input:  `{foo123:=3; bar:=5; return foo123+bar}`,
		output: 8,
	},
	{
		name:   "multiple variables",
		input:  `{al := 3; bal := 5; baz := 10; return bal + al * 4 - baz}`,
		output: 7,
	},
	{
		name:   "return statement with multiple expressions",
		input:  `{{return 1;} 2; 3; return 0}`,
		output: 1,
	},
	{
		name:   "return statement with multiple expressions 2",
		input:  `{1; {return 2;} 3; return 0}`,
		output: 2,
	},
	{
		name:   "return statement with multiple expressions 3",
		input:  `{1; 2; return 3}`,
		output: 3,
	},
	{
		name:   "nested block",
		input:  `{ {1; {2;} return 3;} }`,
		output: 3,
	},
	{
		name:   "return statement with multiple semicolons",
		input:  `{ ;;; return 5; }`,
		output: 5,
	},
	{
		name:   "if statement false without else",
		input:  `{ if false {return 2} return 3 }`,
		output: 3,
	},
	{
		name:   "if statement false with else",
		input:  `{ if false { 1; 2; return 3 } else { return 4 } }`,
		output: 4,
	},
	{
		name:   "if statement true with else",
		input:  `{ if true { 1; 2; return 3 } else { return 4 } }`,
		output: 3,
	},
	{
		name:   "for loop",
		input:  `{ i:=0; j:=0; for i=0; i<=10; i=i+1 { j=i+j; } return j }`,
		output: 55,
	},
	{
		name:   "for loop with empty init and post statements",
		input:  `{ for {return 3;} return 5 }`,
		output: 3,
	},
	{
		name:   "dereference pointer",
		input:  `{ x:=3; return *&x; }`,
		output: 3,
	},
	{
		name:   "dereference pointer to pointer",
		input:  `{ x:=3; y:=&x; z:=&y; return **z; }`,
		output: 3,
	},
	{
		name:   "modify value through pointer",
		input:  `{ x:=3; y:=&x; *y=5; return x; }`,
		output: 5,
	},
	{
		name:   "if expression then value",
		input:  `{ a := if true {1} else {2}; 5; return a }`,
		output: 1,
	},
	{
		name:   "if expression else value",
		input:  `{ a := if false {1} else {2}; 5; return a }`,
		output: 2,
	},
	{
		name:   "block expression",
		input:  `{ a := {true; 42}; return a }`,
		output: 42,
	},
	{
		name:   "reassign",
		input:  `{ a := 1; a = 29; return a }`,
		output: 29,
	},
	{
		name: "function call",
		input: `
			func foo() int {
				return 42
			}
			func main() int {
				return foo()
			}
		`,
		output: 42,
	},
	{
		name: "function call out of order",
		input: `
			func main() int {
				return foo() * 2
			}
			func foo() int {
				return 9
			}
		`,
		output: 18,
	},
	{
		name: "function call with 1 parameter",
		input: `
			func main() int {
				return foo(3) * 3
			}
			func foo(a int) int {
				return a+2
			}
		`,
		output: 15,
	},
	{
		name: "function call with 2 parameters",
		input: `
			func main() int {
				return foo(3, 4) * 3
			}
			func foo(a int, b int) int {
				return a+b
			}
		`,
		output: 21,
	},
	{
		name: "function call with 8 parameters",
		input: `
			func main() int {
				return foo(1, 2, 3, 4, 5, 6, 7, 8)
			}
			func foo(a int, b int, c int, d int, e int, f int, g int, h int) int {
				return a+b+c+d+e+f+g+h
			}
		`,
		output: 36,
	},
}

func TestCodegenWithVirtualMachine(t *testing.T) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := "func main() int " + tt.input
			if strings.Contains(tt.input, "main()") {
				input = tt.input
			}
			file := token.NewFile("test.gos", []byte(input))
			asm := vm.NewAsm()
			errs := compile.Compile(file, asm)

			if len(errs) > 0 {
				for _, err := range errs {
					t.Errorf("Expected no error, but got\n%s", err)
				}
			}

			vm := vm.NewCPU(asm.Program)
			// vm.Trace = true
			actual := vm.Run()
			if actual != tt.output {
				t.Errorf("Expected: %d; but got: %d", tt.output, actual)
			}
		})
	}
}

func TestCodegenNativeAssembly(t *testing.T) {
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := "func main() int " + tt.input
			if strings.Contains(tt.input, "main()") {
				input = tt.input
			}
			file := token.NewFile("test.gos", []byte(input))

			tmp, err := os.CreateTemp("", "gosling_*.s")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmp.Name())

			asm := &aarch64.Assembler{Out: tmp}

			errs := compile.Compile(file, asm)
			if len(errs) > 0 {
				for _, err := range errs {
					t.Errorf("Expected no error, but got %s", err)
				}
				return
			}

			cmd := exec.Command("gcc", "-o", tmp.Name()+".out", tmp.Name())
			err = cmd.Run()
			defer os.Remove(tmp.Name() + ".out")
			if err != nil {
				buf, err := os.ReadFile(tmp.Name())
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(string(buf))
				t.Errorf("Expected no error, but got %s", err)
			}

			cmd = exec.Command(tmp.Name() + ".out")
			_, err = cmd.Output()
			if err != nil {
				if tt.output == 0 {
					t.Errorf("Expected no error, but got %s", err)
					return
				}
				if ee, ok := err.(*exec.ExitError); ok {
					actual := ee.ExitCode()
					if actual != tt.output {
						t.Errorf("Expected: %d; but got: %d", tt.output, actual)
					}
					return
				}
				t.Errorf("Expected no error, but got %s", err)
			}
			if tt.output != 0 && err == nil {
				t.Errorf("Expected error, but got none")
			}
		})
	}
}
