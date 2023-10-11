package hlir_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/compile"
	"github.com/rj45/gosling/hlir"
	"github.com/rj45/gosling/token"
)

var tests = []struct {
	name string
	src  string
	ir   string
}{
	{
		name: "basic function",
		src: `
			func main() int {
				return 42
			}
		`,
		ir: `
			func main() int {
			main.entry0:
				Prologue 0
				r0 = LoadInt 42
				Jump main.epilogue0
			main.epilogue0:
				Epilogue
				Return r0
			}
		`,
	},
	{
		name: "basic function with locals",
		src: `
			func main() int {
				x := 42
				return x
			}
		`,
		ir: `
			func main() int {
			main.entry0:
				Prologue 1
				r0 = LocalAddr 0
				Push r0
				r0 = LoadInt 42
				r1 = Pop
				Store r0, r1
				r0 = LoadLocal 0
				Jump main.epilogue0
			main.epilogue0:
				Epilogue
				Return r0
			}
		`,
	},
	{
		name: "basic function with locals and args",
		src: `
			func main(x int) int {
				return x
			}
		`,
		ir: `
			func main(r0 int) int {
			main.entry0:
				Prologue 1
				StoreLocal r0, 0
				r0 = LoadLocal 0
				Jump main.epilogue0
			main.epilogue0:
				Epilogue
				Return r0
			}
		`,
	},
	{
		name: "function with add, sub, mul, div",
		src: `
			func main() int {
				return 1 + 2 - 3 * 4 / 5
			}
		`,
		ir: `
			func main() int {
			main.entry0:
				Prologue 0
				r0 = LoadInt 1
				Push r0
				r0 = LoadInt 2
				r1 = Pop
				r0 = Add r1, r0
				Push r0
				r0 = LoadInt 3
				Push r0
				r0 = LoadInt 4
				r1 = Pop
				r0 = Mul r1, r0
				Push r0
				r0 = LoadInt 5
				r1 = Pop
				r0 = Div r1, r0
				r1 = Pop
				r0 = Sub r1, r0
				Jump main.epilogue0
			main.epilogue0:
				Epilogue
				Return r0
			}
		`,
	},
	{
		name: "function with neg",
		src: `
			func main() int {
				return -42
			}
		`,
		ir: `
			func main() int {
			main.entry0:
				Prologue 0
				r0 = LoadInt 42
				r0 = Neg r0
				Jump main.epilogue0
			main.epilogue0:
				Epilogue
				Return r0
			}
		`,
	},
	{
		name: "function with eq, ne, lt, gt, le, ge",
		src: `
			func main() int {
				a := 1 < 2
				b := 1 > 2
				c := 1 <= 2
				d := 1 >= 2
				e := a == b
				f := c != d
				if e == f {
					return 1
				}
				return 0
			}
		`,
		ir: `
			func main() int {
            main.entry0:
            	Prologue 6
            	r0 = LocalAddr 0
            	Push r0
            	r0 = LoadInt 1
            	Push r0
            	r0 = LoadInt 2
            	r1 = Pop
            	r0 = Lt r1, r0
            	r1 = Pop
            	Store r0, r1
            	r0 = LocalAddr 1
            	Push r0
            	r0 = LoadInt 1
            	Push r0
            	r0 = LoadInt 2
            	r1 = Pop
            	r0 = Gt r1, r0
            	r1 = Pop
            	Store r0, r1
            	r0 = LocalAddr 2
            	Push r0
            	r0 = LoadInt 1
            	Push r0
            	r0 = LoadInt 2
            	r1 = Pop
            	r0 = Le r1, r0
            	r1 = Pop
            	Store r0, r1
            	r0 = LocalAddr 3
            	Push r0
            	r0 = LoadInt 1
            	Push r0
            	r0 = LoadInt 2
            	r1 = Pop
            	r0 = Ge r1, r0
            	r1 = Pop
            	Store r0, r1
            	r0 = LocalAddr 4
            	Push r0
            	r0 = LoadLocal 0
            	Push r0
            	r0 = LoadLocal 1
            	r1 = Pop
            	r0 = Eq r1, r0
            	r1 = Pop
            	Store r0, r1
            	r0 = LocalAddr 5
            	Push r0
            	r0 = LoadLocal 2
            	Push r0
            	r0 = LoadLocal 3
            	r1 = Pop
            	r0 = Ne r1, r0
            	r1 = Pop
            	Store r0, r1
            	r0 = LoadLocal 4
            	Push r0
            	r0 = LoadLocal 5
            	r1 = Pop
            	r0 = Eq r1, r0
            	If r0, then0, endif0
			then0:
				r0 = LoadInt 1
				Jump main.epilogue0
            endif0:
            	r0 = LoadInt 0
            	Jump main.epilogue0
            main.epilogue0:
            	Epilogue
            	Return r0
            }
		`,
	},
	{
		name: "function with call",
		src: `
			func foo() {}
			func main() int {
				foo()
				return 0
			}
		`,
		ir: `
			func foo() {
			foo.entry0:
				Prologue 0
				Jump foo.epilogue0
			foo.epilogue0:
				Epilogue
				Return
			}

			func main() int {
			main.entry0:
				Prologue 0
				Call foo
				r0 = LoadInt 0
				Jump main.epilogue0
			main.epilogue0:
				Epilogue
				Return r0
			}
		`,
	},
}

func TestAssembler(t *testing.T) {
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			file := token.NewFile("test.gos", []byte(test.src))
			asm := hlir.NewBuilder(file)
			compile.Compile(file, asm)
			actual := asm.Program.Dump()

			if trim(actual) != trim(test.ir) {
				t.Errorf("expected:\n%s\nactual:\n%s", test.ir, actual)
			}
		})
	}
}

func trim(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	return strings.Join(lines, "\n")
}
