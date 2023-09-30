package semantics_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/types"
)

func TestTypeCheckingFuncs(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
		err      string
	}{
		{
			name:     "return needs int",
			src:      "func main() int { return true }",
			expected: "",
			err:      "cannot return bool",
		},
		{
			name:     "invalid return",
			src:      "func main() int { return }",
			expected: "",
			err:      "invalid return statement",
		},
		{
			name:     "function declaration",
			src:      "func foo() int { return 42 }",
			expected: "func() int",
			err:      "",
		},
		{
			name:     "function declaration",
			src:      "func foo() {}",
			expected: "func()",
			err:      "",
		},
		{
			name:     "function redeclaration",
			src:      "func foo() {} func foo() {}",
			expected: "",
			err:      "cannot redefine function foo",
		},
		{
			name:     "bad function return type",
			src:      "func foo() bar {}",
			expected: "",
			err:      "undefined name bar",
		},
		{
			name:     "cannot return value from void function",
			src:      "func foo() {return 1}",
			expected: "",
			err:      "cannot return value from void function",
		},
		{
			name:     "return wrong type",
			src:      "func foo() int {return 1} func bar() bool { return foo() }",
			expected: "",
			err:      "cannot return int from function returning bool",
		},
		{
			name:     "function call to undefined function",
			src:      "func foo() { bar() }",
			expected: "",
			err:      "cannot call undefined function bar",
		},
		{
			name:     "function call with wrong type argument",
			src:      "func foo(a int) {} func bar() { foo(true) }",
			expected: "",
			err:      "wrong type for argument: expected int, got bool",
		},
		{
			name:     "function call with wrong number of arguments",
			src:      "func foo(a int) {} func bar() { foo(1, 2) }",
			expected: "",
			err:      "wrong number of arguments to foo: expected 1, got 2",
		},
		{
			name:     "return wrong type via parameter",
			src:      "func foo(a bool) int { return a }",
			expected: "",
			err:      "cannot return bool from function returning int",
		},
		{
			name:     "function decl",
			src:      "func foo(a bool) int { return 0 }",
			expected: "func(bool) int",
			err:      "",
		},
		{
			name:     "function post declaration",
			src:      "func main() int { return foo(1) } func foo(a int) int { return a }",
			expected: "func() int",
			err:      "",
		},
		{
			name:     "function missing return in function body",
			src:      "func main() int { }",
			expected: "",
			err:      "missing return statement in function main",
		},
		{
			name:     "function return must be last statement in block",
			src:      "func main() int { {return 0; 1+1} return 1 } ",
			expected: "",
			err:      "return must be last statement in block",
		},
		{
			name:     "function return can be in nested block",
			src:      "func main() int { {return 0} } ",
			expected: "func() int",
			err:      "",
		},
		{
			name:     "function return can be in if / else",
			src:      "func main() int { if(true) { return 0 } else { return 1 } }",
			expected: "func() int",
			err:      "",
		},
		{
			name:     "function missing return if only one branch has return",
			src:      "func main() int { if(true) { 1 } else { return 1 } }",
			expected: "",
			err:      "missing return statement in function main",
		},
		{
			name:     "function can return from for loop",
			src:      "func main() int { for { return 1 } }",
			expected: "func() int",
			err:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, uni, node, errs := parse(t, tt.src)

			if errs != nil {
				for _, err := range errs {
					if tt.err == "" {
						t.Errorf("Expected no error, but got %s", errs)
					} else if !strings.Contains(err.Error(), tt.err) {
						t.Errorf("Expected error to contain %q, but got %q", tt.err, err)
					}
				}
				return
			}

			// first child of the root decl list
			node = a.Child(node, 0)

			if errs != nil {
				if tt.err == "" {
					t.Errorf("Expected no error, but got %s", errs)
				} else if !strings.Contains(errs[0].Error(), tt.err) {
					t.Errorf("Expected error to contain %q, but got %q", tt.err, errs[0])
				}

				return
			}

			if errs == nil && tt.err != "" {
				t.Errorf("Expected error %q, but got none", tt.err)
				return
			}

			actual := a.Type(node)

			if actual == types.None && tt.expected != "" {
				t.Errorf("Expected type %q, but got none", tt.expected)
			} else if uni.StringOf(actual) != tt.expected {
				t.Errorf("Expected: %s\nBut got: %s", tt.expected, uni.StringOf(actual))
			}
		})
	}
}
