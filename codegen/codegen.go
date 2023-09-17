package codegen

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

type Assembly interface {
	WordSize() int

	Prologue(int)
	Epilogue()

	Push()
	Pop()
	LoadLocal(int)
	StoreLocal(int)

	LoadInt(string)

	Add()
	Sub()
	Mul()
	Div()

	Neg()

	Eq()
	Ne()
	Lt()
	Le()
	Gt()
	Ge()

	JumpToEpilogue()
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
	g.asm.Prologue(g.ast.StackSize())

	g.genStmtList(g.ast.Root())
	g.asm.Epilogue()
}

func (g *CodeGen) genStmtList(node ast.NodeID) {
	children := g.ast.Children(node)
	for _, child := range children {
		g.genStmt(child)
	}
}

func (g *CodeGen) genStmt(node ast.NodeID) {
	switch g.ast.Kind(node) {
	case ast.ExprStmt:
		g.genExpr(g.ast.Child(node, ast.ExprStmtExpr))
	case ast.AssignStmt:
		g.genAssignStmt(node)
	case ast.ReturnStmt:
		g.genReturnStmt(node)
	default:
		panic("unknown stmt kind")
	}
}

func (g *CodeGen) genAssignStmt(node ast.NodeID) {
	g.genExpr(g.ast.Child(node, ast.AssignStmtRHS))
	g.asm.StoreLocal(g.localOffset(g.ast.Child(node, ast.AssignStmtLHS)))
}

func (g *CodeGen) genReturnStmt(node ast.NodeID) {
	for _, child := range g.ast.Children(node) {
		g.genExpr(child)
	}
	g.asm.JumpToEpilogue()
}

func (g *CodeGen) localOffset(node ast.NodeID) int {
	switch g.ast.Kind(node) {
	case ast.Name:
		sym := g.ast.SymbolOf(node)
		return -(sym.Offset * g.asm.WordSize())
	default:
		panic("unknown addr kind")
	}
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
		case token.Eq:
			g.asm.Eq()
		case token.Ne:
			g.asm.Ne()
		case token.Lt:
			g.asm.Lt()
		case token.Le:
			g.asm.Le()
		case token.Gt:
			g.asm.Gt()
		case token.Ge:
			g.asm.Ge()
		}
	case ast.UnaryExpr:
		g.genExpr(g.ast.Child(node, ast.UnaryExprExpr))
		g.asm.Neg()
	case ast.Literal:
		g.asm.LoadInt(g.ast.NodeString(node))
	case ast.Name:
		g.asm.LoadLocal(g.localOffset(node))
	default:
		panic("unknown expr kind")
	}
}
