package parser_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
)

func TestParsePrimaryExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"42", `Literal("42")`},
		{"foo", `Name("foo")`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte("{" + tt.src + "}"))
		a, err := parser.Parse()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		root := a.Root()
		if a.Kind(root) != ast.StmtList {
			t.Errorf("Expected StmtList, but got %s", a.Kind(root))
		}

		stmt := a.Child(root, 0)
		expr := a.Child(stmt, ast.ExprStmtExpr)

		if trim(a.StringOf(expr)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(expr))
		}
	}
}

func TestParseAddExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"1+2", `BinaryExpr("+", Literal("1"), Literal("2"))`},
		{"1-2", `BinaryExpr("-", Literal("1"), Literal("2"))`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte("{" + tt.src + "}"))
		a, err := parser.Parse()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		root := a.Root()
		if a.Kind(root) != ast.StmtList {
			t.Errorf("Expected StmtList, but got %s", a.Kind(root))
		}

		stmt := a.Child(root, 0)
		expr := a.Child(stmt, ast.ExprStmtExpr)

		if trim(a.StringOf(expr)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(expr))
		}
	}
}

func TestParseMulExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"1+2*3", `
			BinaryExpr("+",
				Literal("1"),
				BinaryExpr("*", Literal("2"), Literal("3")),
			)
		`},
		{"(1+2)/3", `
			BinaryExpr("/",
				BinaryExpr("+", Literal("1"), Literal("2")),
				Literal("3"),
			)
		`},
		{"1+2*3+4", `
			BinaryExpr("+",
				BinaryExpr("+",
					Literal("1"),
					BinaryExpr("*", Literal("2"), Literal("3")),
				),
				Literal("4"),
			)
		`},
		{"1+2*3-4/5", `
			BinaryExpr("-",
				BinaryExpr("+",
					Literal("1"),
					BinaryExpr("*", Literal("2"), Literal("3")),
				),
				BinaryExpr("/", Literal("4"), Literal("5")),
			)
		`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte("{" + tt.src + "}"))
		a, err := parser.Parse()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		root := a.Root()
		if a.Kind(root) != ast.StmtList {
			t.Errorf("Expected StmtList, but got %s", a.Kind(root))
		}

		stmt := a.Child(root, 0)
		expr := a.Child(stmt, ast.ExprStmtExpr)

		if trim(a.StringOf(expr)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(expr))
		}
	}
}

func TestParseUrnaryExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"-5", `UnaryExpr("-", Literal("5"))`},
		{"-5+6", `
			BinaryExpr("+",
				UnaryExpr("-", Literal("5")),
				Literal("6"),
			)
		`},
		{"-5*6", `
			BinaryExpr("*",
				UnaryExpr("-", Literal("5")),
				Literal("6"),
			)
		`},
		{"5-+6", `BinaryExpr("-", Literal("5"), Literal("6"))`},
		{"5-+-+6", `
			BinaryExpr("-",
				Literal("5"),
				UnaryExpr("-", Literal("6")),
			)
		`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte("{" + tt.src + "}"))
		a, err := parser.Parse()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		root := a.Root()
		if a.Kind(root) != ast.StmtList {
			t.Errorf("Expected StmtList, but got %s", a.Kind(root))
		}

		stmt := a.Child(root, 0)
		expr := a.Child(stmt, ast.ExprStmtExpr)

		if trim(a.StringOf(expr)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(expr))
		}
	}
}

func TestParseCompareExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"1==2", `BinaryExpr("==", Literal("1"), Literal("2"))`},
		{"1!=2", `BinaryExpr("!=", Literal("1"), Literal("2"))`},
		{"1<2", `BinaryExpr("<", Literal("1"), Literal("2"))`},
		{"1<=2", `BinaryExpr("<=", Literal("1"), Literal("2"))`},
		{"1>2", `BinaryExpr(">", Literal("1"), Literal("2"))`},
		{"1>=2", `BinaryExpr(">=", Literal("1"), Literal("2"))`},
	}

	for _, tt := range tests {
		parser := parser.New([]byte("{" + tt.src + "}"))
		a, err := parser.Parse()
		if err != nil {
			t.Errorf("Expected no error, but got %s", err)
		}

		root := a.Root()
		if a.Kind(root) != ast.StmtList {
			t.Errorf("Expected StmtList, but got %s", a.Kind(root))
		}

		stmt := a.Child(root, 0)
		expr := a.Child(stmt, ast.ExprStmtExpr)

		if trim(a.StringOf(expr)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(expr))
		}
	}
}

func trim(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, l := range lines {
		lines[i] = strings.TrimSpace(l)
	}
	return strings.Join(lines, "\n")
}
