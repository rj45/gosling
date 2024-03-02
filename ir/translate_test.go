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
				$return = Return $start, 42
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
				$x = Copy 42
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
				$add = Add 1, 2
				$mul = Mul 3, 4
				$div = Div $mul, 5
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
				$neg = Neg 42
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

				$if = If $start, true

				$then = Then $if
				$return1 = Return $then, 1

				$region = Region $start, $return1
				$return2 = Return $region, 0
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

				$if = If $start, true

				$then = Then $if
				$return1 = Return $then, 1

				$else = Else $if
				$return2 = Return $else, 0

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
				$x1 = Copy 0

				$if = If $start, true

				$then = Then $if
				$x2 = Copy 1

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
