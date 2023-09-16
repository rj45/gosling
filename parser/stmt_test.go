package parser_test

import (
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
)

func TestParseExprStmt(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"42", `ExprStmt(Literal("42"))`},
		{"1+2", `ExprStmt(
        	BinaryExpr("+", Literal("1"), Literal("2")),
        )`},
		{"1-2", `ExprStmt(
        	BinaryExpr("-", Literal("1"), Literal("2")),
        )`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte(tt.src))
		a := parser.Parse()

		root := a.Root()
		if a.Kind(root) != ast.StmtList {
			t.Errorf("Expected StmtList, but got %s", a.Kind(root))
		}

		stmt := a.Child(root, 0)

		if trim(a.StringOf(stmt)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(stmt))
		}
	}
}

func TestParseStmtList(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"1+2;3+4", `StmtList(
			ExprStmt(
				BinaryExpr("+", Literal("1"), Literal("2")),
			),
			ExprStmt(
				BinaryExpr("+", Literal("3"), Literal("4")),
			),
		)`},
		{"1+2\n3+4", `StmtList(
			ExprStmt(
				BinaryExpr("+", Literal("1"), Literal("2")),
			),
			ExprStmt(
				BinaryExpr("+", Literal("3"), Literal("4")),
			),
		)`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte(tt.src))
		a := parser.Parse()

		root := a.Root()
		if a.Kind(root) != ast.StmtList {
			t.Errorf("Expected StmtList, but got %s", a.Kind(root))
		}

		if trim(a.StringOf(root)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(root))
		}
	}
}
