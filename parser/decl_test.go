package parser_test

import "testing"

func TestParseFuncDecl(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"func main() int {}", `FuncDecl(
			Name("main"),
			nil,
			Name("int"),
			StmtList(),
		)`},
		{"func main() int { return 42 }", `FuncDecl(
			Name("main"),
			nil,
			Name("int"),
			StmtList(
				ReturnStmt(Literal("42")),
			),
		)`},
		{"func main() int { return 1+2 }", `FuncDecl(
			Name("main"),
			nil,
			Name("int"),
			StmtList(
				ReturnStmt(
					BinaryExpr("+", Literal("1"), Literal("2")),
				),
			),
		)`},
	}

	for _, tt := range tests {
		a, decllist, errs := parse(t, tt.src)
		if len(errs) > 0 {
			t.Errorf("Expected no error, but got %s", errs)
		}

		decl := a.Child(decllist, 0)

		if trim(a.StringOf(decl)) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, a.StringOf(decl))
		}
	}
}
