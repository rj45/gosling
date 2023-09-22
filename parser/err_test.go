package parser_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
)

func TestExprParseError(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"{1+2*", `expected expression`},
		{"{1+2)", `expected newline or ';' or '}'`},
		{"{1+2;", `expected '}'`},
		{"{1foo}", `illegal token "1foo"`},
	}

	for _, tt := range tests {
		parser := parser.New(ast.NewFile("test.gos", []byte(tt.src)))
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
		{"{for i = 9 {}}", `expected for condition to be expression statement`},
		{"{9 = 45}", `expected name or deref on the left side of the assignment`},
	}

	for _, tt := range tests {
		parser := parser.New(ast.NewFile("test.gos", []byte(tt.src)))
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
