package ir

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
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
				$start = Start
				$const = Const 42
				$return = Return $start, $const
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
				$start = Start
				$const = Const 42
				$x = Copy $const
				$return = Return $start, $x
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
				$start = Start
				$const1 = Const 1
				$const2 = Const 2
				$add = Add $const1, $const2
				$const3 = Const 3
				$const4 = Const 4
				$mul = Mul $const3, $const4
				$const5 = Const 5
				$div = Div $mul, $const5
				$sub = Sub $add, $div
				$return = Return $start, $sub
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
				$start = Start
				$const = Const 42
				$neg = Neg $const
				$return = Return $start, $neg
			}
		`,
	},
	{
		name: "function with simple if, no phi",
		src: `
			func main() int {
				if true {
					return 1
				}
				return 0
			}
		`,
		ir: `
			func main() int {
				$start = Start
				$const1 = Const true

				$if = If $start, $const1

				$then = Then $if
				$const2 = Const 1
				$return1 = Return $then, $const2

				$region = Region $start, $return1
				$const3 = Const 0
				$return2 = Return $region, $const3
			}
		`,
	},
	{
		name: "function with simple if-else, no phi",
		src: `
			func main() int {
				if true {
					return 1
				} else {
					return 0
				}
			}
		`,
		ir: `
			func main() int {
				$start = Start
				$const1 = Const true

				$if = If $start, $const1

				$then = Then $if
				$const2 = Const 1
				$return1 = Return $then, $const2

				$else = Else $if
				$const3 = Const 0
				$return2 = Return $else, $const3

				$region = Region $return1, $return2
			}
		`,
	},
	{
		name: "function with simple if with phi",
		src: `
			func main() int {
				x := 0
				if true {
					x = 1
				}
				return x
			}
		`,
		ir: `
			func main() int {
				$start = Start
				$const1 = Const 0
				$x1 = Copy $const1
				$const2 = Const true

				$if = If $start, $const2

				$then = Then $if
				$const3 = Const 1
				$x2 = Copy $const3

				$region = Region $start, $then
				$x3 = Phi $region, $x1, $x2
				$return = Return $region, $x3
			}
		`,
	},
}

func TestTranslator(t *testing.T) {
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			file := ast.NewFile("test.gos", []byte(test.src))

			parser := parser.New(file)
			ast, errs := parser.Parse()
			for _, err := range errs {
				t.Error(err)
			}

			translator := NewSoNTranslator(ast)
			program := translator.Translate()

			actual := program.Dump()

			if translator.errs != nil {
				for _, err := range translator.errs {
					t.Error(err)
				}
			}

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
