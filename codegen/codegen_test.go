package codegen_test

import (
	"fmt"
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/codegen"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/vm"
)

func execute(a *ast.AST) int {
	asm := vm.NewAsm()
	codegen := codegen.New(a, asm)
	codegen.Generate()
	vm := vm.NewCPU(asm.Program)
	return vm.Run()
}

func assert(t *testing.T, result int, prog string) {
	t.Run(fmt.Sprintf("assert %d %q", result, prog), func(t *testing.T) {
		parser := parser.New([]byte(prog))
		a, err := parser.Parse()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}
		actual := execute(a)
		if actual != result {
			t.Errorf("Expected: %d; but got: %d", result, actual)
		}
	})
}

func TestCodegen(t *testing.T) {
	assert(t, 0, `{return 0}`)
	assert(t, 42, `{return 42}`)
	assert(t, 1, `{return 1+2-2}`)
	assert(t, 21, `{return 5 + 20 - 4}`)
	assert(t, 41, `{return  12 + 34 - 5 }`)
	assert(t, 47, `{return 5+6*7}`)
	assert(t, 15, `{return 5*(9-6)}`)
	assert(t, 4, `{return (3+5)/2}`)
	assert(t, 10, `{return -10+20}`)
	assert(t, 10, `{return - -10}`)
	assert(t, 10, `{return - - +10}`)

	assert(t, 0, `{return 0==1}`)
	assert(t, 1, `{return 42==42}`)
	assert(t, 1, `{return 0!=1}`)
	assert(t, 0, `{return 42!=42}`)

	assert(t, 1, `{return 0<1}`)
	assert(t, 0, `{return 1<1}`)
	assert(t, 0, `{return 2<1}`)
	assert(t, 1, `{return 0<=1}`)
	assert(t, 1, `{return 1<=1}`)
	assert(t, 0, `{return 2<=1}`)
	assert(t, 1, `{return 1>0}`)
	assert(t, 0, `{return 1>1}`)
	assert(t, 0, `{return 1>2}`)
	assert(t, 1, `{return 1>=0}`)
	assert(t, 1, `{return 1>=1}`)

	assert(t, 3, `{a=3; return a}`)
	assert(t, 8, `{a=3; z=5; return a+z}`)

	assert(t, 3, `{foo=3; return foo}`)
	assert(t, 8, `{foo123=3; bar=5; return foo123+bar}`)
	assert(t, 7, `{al = 3; bal = 5; baz = 10; return bal + al * 4 - baz}`)
	assert(t, 1, `{return 1; 2; 3}`)
	assert(t, 2, `{1; return 2; 3}`)
	assert(t, 3, `{1; 2; return 3}`)
	assert(t, 3, `{ {1; {2;} return 3;} }`)
	assert(t, 5, `{ ;;; return 5; }`)
	assert(t, 3, `{ if (0) {return 2} return 3 }`)
	assert(t, 3, `{ if (1-1) {return 2} return 3 }`)
	assert(t, 2, `{ if (1) {return 2} return 3 }`)
	assert(t, 2, `{ if (2-1) {return 2} return 3 }`)
	assert(t, 4, `{ if (0) { 1; 2; return 3 } else { return 4 } }`)
	assert(t, 3, `{ if (1) { 1; 2; return 3 } else { return 4 } }`)
	assert(t, 55, `{ i=0; j=0; for i=0; i<=10; i=i+1 { j=i+j; } return j }`)
	assert(t, 3, `{ for {return 3;} return 5 }`)
	assert(t, 3, `{ x=3; return *&x; }`)
	assert(t, 3, `{ x=3; y=&x; z=&y; return **z; }`)
	assert(t, 5, `{ x=3; y=&x; *y=5; return x; }`)
}
