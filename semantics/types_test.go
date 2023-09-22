package semantics_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/semantics"
)

func ParseStmt(src string) (*ast.AST, ast.NodeID, error) {
	parser := parser.New([]byte("{" + src + "}"))
	a, err := parser.Parse()
	if err != nil {
		return nil, ast.InvalidNode, err
	}
	errs := semantics.NewTypeChecker(a).Check(a.Root())
	if len(errs) > 0 {
		return nil, ast.InvalidNode, errs[0]
	}
	node := a.Root()
	node = a.Child(node, a.NumChildren(node)-1)
	return a, node, nil
}

func TestTypeChecking(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
		err      string
	}{
		{
			name:     "int literal is untyped int",
			src:      "42",
			expected: "untyped int",
			err:      "",
		},
		{
			name:     "true is bool",
			src:      "true",
			expected: "bool",
			err:      "",
		},
		{
			name:     "false is bool",
			src:      "false",
			expected: "bool",
			err:      "",
		},
		{
			name:     "adding two int literals is untyped int",
			src:      "1+2",
			expected: "untyped int",
			err:      "",
		},
		{
			name:     "less than is bool",
			src:      "1<2",
			expected: "bool",
			err:      "",
		},
		{
			name:     "equal is bool",
			src:      "1==2",
			expected: "bool",
			err:      "",
		},
		{
			name:     "unary minus is untyped int",
			src:      "-1",
			expected: "untyped int",
			err:      "",
		},
		{
			name:     "if condition needs bool",
			src:      "if 1 {}",
			expected: "",
			err:      "if condition must be bool",
		},
		{
			name:     "for condition needs bool",
			src:      "for 1 {}",
			expected: "",
			err:      "for condition must be bool",
		},
		{
			name:     "for ;; condition needs bool",
			src:      "for ;1; {}",
			expected: "",
			err:      "for condition must be bool",
		},
		{
			name:     "return needs int",
			src:      "return true",
			expected: "",
			err:      "cannot return bool",
		},
		{
			name:     "assign untyped int converted to int",
			src:      "a = 1",
			expected: "int",
			err:      "",
		},
		{
			name:     "reassign true to false",
			src:      "a = true; a = false",
			expected: "bool",
			err:      "",
		},
		{
			name:     "reassign int to bool",
			src:      "a = 1; a = true",
			expected: "",
			err:      "cannot assign bool to int",
		},
		{
			name:     "address of int",
			src:      "a = 1; b = &a",
			expected: "*int",
			err:      "",
		},
		{
			name:     "address of addr of int",
			src:      "a = 1; b = &a; c = &b",
			expected: "**int",
			err:      "",
		},
		{
			name:     "deref int",
			src:      "a = 1; b = &a; c = *b",
			expected: "int",
			err:      "",
		},
		{
			name:     "bad assign of deref",
			src:      "a = 1; b = &a; c = &b; c = b",
			expected: "",
			err:      "cannot assign *int to **int",
		},
		{
			name:     "bad deref of int",
			src:      "a = 1; b = *a",
			expected: "",
			err:      "cannot dereference non-pointer type int",
		},
		{
			name:     "bad address of int",
			src:      "&1",
			expected: "",
			err:      "cannot take address of non-name",
		},
		{
			name:     "invalid return",
			src:      "return",
			expected: "",
			err:      "invalid return statement",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, node, err := ParseStmt(tt.src)

			if err != nil {
				if tt.err == "" {
					t.Errorf("Expected no error, but got %s", err)
				} else if !strings.Contains(err.Error(), tt.err) {
					t.Errorf("Expected error to contain %q, but got %q", tt.err, err)
				}
				return
			}

			if err == nil && tt.err != "" {
				t.Errorf("Expected error %q, but got none", tt.err)
				return
			}

			actual := a.Type(node)

			if actual == nil && tt.expected != "" {
				t.Errorf("Expected type %q, but got none", tt.expected)
			} else if actual.String() != tt.expected {
				t.Errorf("Expected: %s\nBut got: %s", tt.expected, actual)
			}
		})
	}
}
