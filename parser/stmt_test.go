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
		a, stmt, errs := parseStmt(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		if trim(a.StringOf(stmt)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(stmt))
		}
	}
}

func TestParseReturnStmt(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"return 42", `ReturnStmt(Literal("42"))`},
		{"return 1+2", `ReturnStmt(
			BinaryExpr("+", Literal("1"), Literal("2")),
		)`},
		{"return", `ReturnStmt()`},
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

func TestParseForStmt(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"for {1}", `ForStmt(
			nil,
			nil,
			nil,
			StmtList(
				ExprStmt(Literal("1")),
			),
		)`},
		{"for 1 {2}", `ForStmt(
			nil,
			ExprStmt(Literal("1")),
			nil,
			StmtList(
				ExprStmt(Literal("2")),
			),
		)`},
		{"for 1; 2 {3}", `ForStmt(
			ExprStmt(Literal("1")),
			ExprStmt(Literal("2")),
			nil,
			StmtList(
				ExprStmt(Literal("3")),
			),
		)`},
		{"for 1;2;3 {4}", `ForStmt(
			ExprStmt(Literal("1")),
			ExprStmt(Literal("2")),
			ExprStmt(Literal("3")),
			StmtList(
				ExprStmt(Literal("4")),
			),
		)`},
		{"for 1;2;3 {for 4;5;6 {7}}", `ForStmt(
			ExprStmt(Literal("1")),
			ExprStmt(Literal("2")),
			ExprStmt(Literal("3")),
			StmtList(
				ForStmt(
					ExprStmt(Literal("4")),
					ExprStmt(Literal("5")),
					ExprStmt(Literal("6")),
					StmtList(
						ExprStmt(Literal("7")),
					),
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

func TestParseAssignStmt(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"foo=42", `AssignStmt("=", Name("foo"), Literal("42"))`},
		{"foo=1+2", `AssignStmt("=",
			Name("foo"),
			BinaryExpr("+", Literal("1"), Literal("2")),
		)`},
		{"foo=1-2", `AssignStmt("=",
			Name("foo"),
			BinaryExpr("-", Literal("1"), Literal("2")),
		)`},
		{"foo:=42", `AssignStmt(":=", Name("foo"), Literal("42"))`},
		{"foo:=1*2", `AssignStmt(":=",
			Name("foo"),
			BinaryExpr("*", Literal("1"), Literal("2")),
		)`},
		{"*foo = 42", `AssignStmt("=",
			DerefExpr(Name("foo")),
			Literal("42"),
		)`},
		{"**foo = 42", `AssignStmt("=",
        	DerefExpr(
        		DerefExpr(Name("foo")),
        	),
        	Literal("42"),
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

func parse(t *testing.T, src string) (*ast.AST, ast.NodeID, []error) {
	t.Helper()

	parser := parser.New(ast.NewFile("test.gos", []byte(src)))
	a, errs := parser.Parse()
	if len(errs) > 0 {
		return nil, ast.InvalidNode, errs
	}

	root := a.Root()
	if a.Kind(root) != ast.DeclList {
		t.Errorf("Expected DeclList, but got %s", a.Kind(root))
	}

	return a, root, nil
}

func parseStmt(t *testing.T, src string) (*ast.AST, ast.NodeID, []error) {
	t.Helper()

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

	stmt := a.Child(body, 0)
	return a, stmt, nil
}
