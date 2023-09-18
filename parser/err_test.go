package parser_test

import (
	"strings"
	"testing"

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
		parser := parser.New([]byte(tt.src))
		a, err := parser.Parse()

		if err == nil {
			t.Errorf("Expected error, but got %s", a)
		}

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
		{"{9 = 45}", `expected name on the left side of the assignment`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte(tt.src))
		a, err := parser.Parse()

		if err == nil {
			t.Errorf("Expected error, but got %s", a)
		}

		if !strings.Contains(err.Error(), tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, err.Error())
		}
	}
}
