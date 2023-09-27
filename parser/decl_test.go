package parser_test

import "testing"

func TestParseFuncDecl(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"func main() int {}", `FuncDecl(
			Name("main"),
			FieldList(),
			Name("int"),
			StmtList(),
		)`},
		{"func main() int { return 42 }", `FuncDecl(
			Name("main"),
			FieldList(),
			Name("int"),
			StmtList(
				ReturnStmt(Literal("42")),
			),
		)`},
		{"func main() int { return 1+2 }", `FuncDecl(
			Name("main"),
			FieldList(),
			Name("int"),
			StmtList(
				ReturnStmt(
					BinaryExpr("+", Literal("1"), Literal("2")),
				),
			),
		)`},
		{"func foo() {}\n\nfunc main() int { return 1 }", `FuncDecl(
        	Name("foo"),
        	FieldList(),
        	nil,
        	StmtList(),
        )`},
		{"func foo(a int) {}", `FuncDecl(
			Name("foo"),
			FieldList(
				Field(Name("a"), Name("int")),
			),
			nil,
			StmtList(),
		)`},
		{"func foo(a int, b bool) int {}", `FuncDecl(
			Name("foo"),
			FieldList(
				Field(Name("a"), Name("int")),
				Field(Name("b"), Name("bool")),
			),
			Name("int"),
			StmtList(),
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
