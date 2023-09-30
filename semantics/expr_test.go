package semantics_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/types"
)

func TestTypeCheckingExprs(t *testing.T) {
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
			name:     "address of int",
			src:      "a := 1; b := &a",
			expected: "*int",
			err:      "",
		},
		{
			name:     "address of addr of int",
			src:      "a := 1; b := &a; c := &b",
			expected: "**int",
			err:      "",
		},
		{
			name:     "deref int",
			src:      "a := 1; b := &a; c := *b",
			expected: "int",
			err:      "",
		},
		{
			name:     "bad assign of deref",
			src:      "a := 1; b := &a; c := &b; c = b",
			expected: "",
			err:      "cannot assign *int to **int",
		},
		{
			name:     "bad deref of int",
			src:      "a := 1; b := *a",
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
			name:     "if expression",
			src:      "a := if true {1} else {2}; a",
			expected: "int",
			err:      "",
		},
		{
			name:     "if expression with mismatched types",
			src:      "a := if true {true} else {2}",
			expected: "",
			err:      "if branches have mismatched types: bool and untyped int",
		},
		{
			name:     "if expression as statement does not get error",
			src:      "if true {true} else {2}",
			expected: "bool",
			err:      "",
		},
		{
			name:     "if expression with untyped int becomes int",
			src:      "a := 2; b := if true {2} else {a}; b",
			expected: "int",
			err:      "",
		},
		{
			name:     "block expression",
			src:      "a := {true; 2}; a",
			expected: "int",
			err:      "",
		},
		{
			name:     "block expression",
			src:      "a := true; a = {true; 2}; a",
			expected: "",
			err:      "cannot assign int to bool",
		},
		{
			name:     "variable use before definition",
			src:      "a = 1",
			expected: "",
			err:      "undefined name a",
		},
		{
			name:     "function call to undefined function",
			src:      "a := bar(); a",
			expected: "",
			err:      "cannot call undefined function bar",
		},
		{
			name:     "function call to non function",
			src:      "a := 1; a()",
			expected: "",
			err:      "cannot call non-function a of type int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, uni, node, errs := parseStmt(t, tt.src)

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
