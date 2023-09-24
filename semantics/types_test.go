package semantics_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/semantics"
)

func parse(t *testing.T, src string) (*ast.AST, ast.NodeID, []error) {
	parser := parser.New(ast.NewFile("test.gos", []byte(src)))
	a, errs := parser.Parse()
	if len(errs) > 0 {
		return nil, ast.InvalidNode, errs
	}

	root := a.Root()
	if a.Kind(root) != ast.DeclList {
		t.Errorf("Expected DeclList, but got %s", a.Kind(root))
	}

	errs = semantics.NewTypeChecker(a).Check(root)
	if errs != nil {
		return nil, ast.InvalidNode, errs
	}

	return a, root, nil
}

func parseStmt(t *testing.T, src string) (*ast.AST, ast.NodeID, []error) {
	a, root, errs := parse(t, "func main() int {"+src+"}")
	if len(errs) > 0 {
		return nil, ast.InvalidNode, errs
	}

	decl := a.Child(root, 0)
	if a.Kind(decl) != ast.FuncDecl {
		t.Errorf("Expected FuncDecl, but got %s", a.Kind(decl))
	}

	body := a.Child(decl, ast.FuncDeclBody)
	if a.Kind(body) != ast.StmtList {
		t.Errorf("Expected StmtList, but got %s", a.Kind(root))
	}

	stmt := a.Child(body, a.NumChildren(body)-1)
	return a, stmt, nil
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
			src:      "a := 1",
			expected: "int",
			err:      "",
		},
		{
			name:     "reassign true to false",
			src:      "a := true; a = false",
			expected: "bool",
			err:      "",
		},
		{
			name:     "reassign int to bool",
			src:      "a := 1; a = true",
			expected: "",
			err:      "cannot assign bool to int",
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
			name:     "invalid return",
			src:      "return",
			expected: "",
			err:      "invalid return statement",
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
			name:     "variable redefinition",
			src:      "a := 1; a := true",
			expected: "",
			err:      "cannot redefine a",
		},
		{
			name:     "function call to undefined function",
			src:      "a := foo(); a",
			expected: "",
			err:      "cannot call undefined function foo",
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
			a, node, errs := parseStmt(t, tt.src)

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

			if actual == nil && tt.expected != "" {
				t.Errorf("Expected type %q, but got none", tt.expected)
			} else if actual.String() != tt.expected {
				t.Errorf("Expected: %s\nBut got: %s", tt.expected, actual)
			}
		})
	}
}

func TestTypeCheckingDecls(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected string
		err      string
	}{
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
			err:      "function foo has invalid return type bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, node, errs := parse(t, tt.src)

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

			if actual == nil && tt.expected != "" {
				t.Errorf("Expected type %q, but got none", tt.expected)
			} else if actual.String() != tt.expected {
				t.Errorf("Expected: %s\nBut got: %s", tt.expected, actual)
			}
		})
	}
}
