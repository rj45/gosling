package parser_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/ast"
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
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if a.Kind(stmt) != ast.ExprStmt {
			t.Errorf("Expected ExprStmt, but got %s", a.Kind(stmt))
		}

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
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if a.Kind(stmt) != ast.ExprStmt {
			t.Errorf("Expected ExprStmt, but got %s", a.Kind(stmt))
		}

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
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if a.Kind(stmt) != ast.ExprStmt {
			t.Errorf("Expected ExprStmt, but got %s", a.Kind(stmt))
		}

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
		{"*foo", `DerefExpr(Name("foo"))`},
		{"**foo", `DerefExpr(
			DerefExpr(Name("foo")),
		)`},
		{"&foo", `AddrExpr(Name("foo"))`},
	}

	for _, tt := range tests {
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if a.Kind(stmt) != ast.ExprStmt {
			t.Errorf("Expected ExprStmt, but got %s", a.Kind(stmt))
		}

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
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if a.Kind(stmt) != ast.ExprStmt {
			t.Errorf("Expected ExprStmt, but got %s", a.Kind(stmt))
		}

		expr := a.Child(stmt, ast.ExprStmtExpr)

		if trim(a.StringOf(expr)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(expr))
		}
	}
}

func TestParseBlock(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"{42}", `StmtList(
			ExprStmt(Literal("42")),
		)`},
		{"{42;12}", `StmtList(
			ExprStmt(Literal("42")),
			ExprStmt(Literal("12")),
		)`},
		{"{42\n12}", `StmtList(
			ExprStmt(Literal("42")),
			ExprStmt(Literal("12")),
		)`},
		{"{1+2}", `StmtList(
			ExprStmt(
				BinaryExpr("+", Literal("1"), Literal("2")),
			),
		)`},
		{"{ {1; {2;} return 3;} }", `StmtList(
			StmtList(
				ExprStmt(Literal("1")),
				StmtList(
					ExprStmt(Literal("2")),
				),
				ReturnStmt(Literal("3")),
			),
		)`},
		{"{ ;;;; }", `StmtList(
			EmptyStmt(),
			EmptyStmt(),
			EmptyStmt(),
			EmptyStmt(),
		)`},
		{"return {42}", `ReturnStmt(
			StmtList(
				ExprStmt(Literal("42")),
			),
		)`},
	}

	for _, tt := range tests {
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if trim(a.StringOf(stmt)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(stmt))
		}
	}
}

func TestParseIfExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"if 1 {2}", `IfExpr(
			Literal("1"),
			StmtList(
				ExprStmt(Literal("2")),
			),
			nil,
		)`},
		{"if 1 {2} else {3}", `IfExpr(
			Literal("1"),
			StmtList(
				ExprStmt(Literal("2")),
			),
			StmtList(
				ExprStmt(Literal("3")),
			),
		)`},
		{"if 1 {2} else {if 3 {4} else {5}}", `IfExpr(
			Literal("1"),
			StmtList(
				ExprStmt(Literal("2")),
			),
			StmtList(
				IfExpr(
					Literal("3"),
					StmtList(
						ExprStmt(Literal("4")),
					),
					StmtList(
						ExprStmt(Literal("5")),
					),
				),
			),
		)`},
		{"a = if true {1} else {2}", `AssignStmt("=",
			Name("a"),
			IfExpr(
				Name("true"),
				StmtList(
					ExprStmt(Literal("1")),
				),
				StmtList(
					ExprStmt(Literal("2")),
				),
			),
		)`},
	}
	for _, tt := range tests {
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if trim(a.StringOf(stmt)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(stmt))
		}
	}
}

func TestParseCallExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"foo()", `CallExpr(
			Name("foo"),
			ExprList(),
        )`},
		{"foo(1)", `CallExpr(
			Name("foo"),
			ExprList(Literal("1")),
        )`},
		{"foo(1, 2)", `CallExpr(
			Name("foo"),
			ExprList(Literal("1"), Literal("2")),
        )`},
		{"foo(1+2, 4*8)", `CallExpr(
			Name("foo"),
			ExprList(
				BinaryExpr("+", Literal("1"), Literal("2")),
				BinaryExpr("*", Literal("4"), Literal("8")),
			),
		)`},
	}
	for _, tt := range tests {
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if a.Kind(stmt) != ast.ExprStmt {
			t.Errorf("Expected ExprStmt, but got %s", a.Kind(stmt))
		}

		expr := a.Child(stmt, ast.ExprStmtExpr)

		if trim(a.StringOf(expr)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(stmt))
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
