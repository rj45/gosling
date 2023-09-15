package parser_test

import (
	"strings"
	"testing"

	"github.com/rj45/gosling/parser"
)

func TestParseExpr(t *testing.T) {
	tests := []struct {
		src      string
		expected string
	}{
		{"42", `Literal("42")`},
		{"1+2", `BinaryExpr("+", Literal("1"), Literal("2"))`},
		{"1-2", `BinaryExpr("-", Literal("1"), Literal("2"))`},
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
		parser := parser.New([]byte(tt.src))
		ast := parser.Parse()

		if trim(ast.String()) != trim(tt.expected) {
			t.Errorf("Expected: %s\nBut got: %s", tt.expected, ast.String())
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
