package parser_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/token"
)

func TestExprParseError(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"func main() int {1+2*", `expected expression`},
		{"func main() int {1+2)", `expected newline or ';' or '}'`},
		{"func main() int {1+2;", `expected '}'`},
		{"func main() int {1foo}", `illegal token "1foo"`},
	}

	for _, tt := range tests {
		parser := parser.New(token.NewFile("test.gos", []byte(tt.src)))
		a, errs := parser.Parse()

		if errs == nil {
			t.Errorf("Expected error, but got %s", a)
		}

		err := errs[0]

		if !strings.Contains(err.Error(), tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, err.Error())
		}
	}
}

func TestStmtParseError(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"func main() int {for i = 9 {}}", `expected for condition to be expression statement`},
		{"func main() int {9 = 45}", `expected name or deref on the left side of the assignment`},
	}

	for _, tt := range tests {
		parser := parser.New(token.NewFile("test.gos", []byte(tt.src)))
		a, errs := parser.Parse()

		if errs == nil {
			t.Errorf("Expected error, but got %s", a)
		}

		err := errs[0]

		if !strings.Contains(err.Error(), tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, err.Error())
		}
	}
}

func TestDeclParseError(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"{", `expected declaration`},
		{"func foo(a) {}", `expected type`},
		{"func(a int) {}", `expected name`},
		{"func foo()\n{}", `expected '{'`}, // bracket must be on same line (like in Go)
		{"func foo() {", `expected '}'`},
		{"func foo {}", `expected '('`},
	}

	for _, tt := range tests {
		parser := parser.New(token.NewFile("test.gos", []byte(tt.src)))
		a, errs := parser.Parse()

		if errs == nil {
			t.Errorf("Expected error, but got %s", a)
		}

		err := errs[0]

		if !strings.Contains(err.Error(), tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, err.Error())
		}
	}
}
