package codegen

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

type Assembly interface {
	Prologue()
	Epilogue()

	Push()
	Pop()

	LoadInt(string)

	Add()
	Sub()
	Mul()
	Div()

	Neg()
}

type CodeGen struct {
	ast *ast.AST
	asm Assembly
}

func New(ast *ast.AST, asm Assembly) *CodeGen {
	return &CodeGen{
		ast: ast,
		asm: asm,
	}
}

func (g *CodeGen) Generate() {
	g.asm.Prologue()
	g.genExpr(g.ast.Root())
	g.asm.Epilogue()
}

func (g *CodeGen) genExpr(node ast.NodeID) {
	switch g.ast.Kind(node) {
	case ast.BinaryExpr:
		g.genExpr(g.ast.Child(node, ast.BinaryExprLHS))
		g.asm.Push()
		g.genExpr(g.ast.Child(node, ast.BinaryExprRHS))
		g.asm.Pop()

		switch g.ast.Token(node).Kind() {
		case token.Add:
			g.asm.Add()
		case token.Sub:
			g.asm.Sub()
		case token.Mul:
			g.asm.Mul()
		case token.Div:
			g.asm.Div()
		}
	case ast.UnaryExpr:
		g.genExpr(g.ast.Child(node, ast.UnaryExprExpr))
		g.asm.Neg()
	case ast.Literal:
		g.asm.LoadInt(g.ast.NodeString(node))
	default:
		panic("unknown expr kind")
	}
}
