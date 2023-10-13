package codegen

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

func (g *CodeGen) genIfExpr(node ast.NodeID) {
	cond := g.ast.Child(node, ast.IfExprCond)
	then := g.ast.Child(node, ast.IfExprThen)
	els := g.ast.Child(node, ast.IfExprElse)

	label := g.label
	g.label++

	g.genExpr(cond)
	if els != ast.InvalidNode {
		g.asm.JumpIfFalse("else", label)
	} else {
		g.asm.JumpIfFalse("endif", label)
	}
	g.asm.Label("then", label)
	g.genStmt(then, false)
	g.asm.Jump("endif", label)
	if els != ast.InvalidNode {
		g.asm.Label("else", label)
		g.genStmt(els, false)
	}
	g.asm.Label("endif", label)
}

func (g *CodeGen) localOffset(node ast.NodeID) int {
	switch g.ast.Kind(node) {
	case ast.Name:
		sym := g.symtab.Lookup(g.ast.NodeString(node))
		return sym.Offset * g.asm.WordSize()
	default:
		panic("unknown addr kind for offset")
	}
}

func (g *CodeGen) genExpr(node ast.NodeID) {
	switch g.ast.Kind(node) {
	case ast.BinaryExpr:
		g.genExpr(g.ast.Child(node, ast.BinaryExprLHS))
		g.asm.Push()
		g.genExpr(g.ast.Child(node, ast.BinaryExprRHS))
		g.asm.Pop(1)

		switch g.ast.Token(node).Kind() {
		case token.Add:
			g.asm.Add()
		case token.Sub:
			g.asm.Sub()
		case token.Star:
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
	case ast.DerefExpr:
		g.genExpr(g.ast.Child(node, ast.DerefExprExpr))
		g.asm.Load()
	case ast.AddrExpr:
		g.genAddr(g.ast.Child(node, ast.AddrExprExpr))
	case ast.IfExpr:
		g.genIfExpr(node)
	case ast.StmtList:
		g.genStmtList(node, false)
	case ast.Literal:
		g.asm.LoadInt(g.ast.NodeString(node))
	case ast.Name:
		sym := g.symtab.Lookup(g.ast.NodeString(node))
		if sym.Const != nil {
			g.genConst(sym.Const)
			return
		}
		g.asm.LoadLocal(g.localOffset(node))
	case ast.CallExpr:
		g.genCallExpr(node)
	default:
		panic("unknown expr kind")
	}
}

func (g *CodeGen) genAddr(node ast.NodeID) {
	switch g.ast.Kind(node) {
	case ast.Name:
		g.asm.LocalAddr(g.localOffset(node))
	case ast.DerefExpr:
		g.genExpr(g.ast.Child(node, ast.DerefExprExpr))
	default:
		panic("unknown addr kind")
	}
}

func (g *CodeGen) genConst(c types.Const) {
	if _, ok := types.Int64Value(c); ok {
		g.asm.LoadInt(c.String())
		return
	}

	if val, ok := types.BoolValue(c); ok {
		if val {
			g.asm.LoadInt("1")
		} else {
			g.asm.LoadInt("0")
		}
		return
	}
	panic("todo: implement const type " + c.String())
}

func (g *CodeGen) genCallExpr(node ast.NodeID) {
	name := g.ast.Child(node, ast.CallExprName)
	argList := g.ast.Child(node, ast.CallExprArgs)
	for _, arg := range g.ast.Children(argList) {
		g.genExpr(arg)
		g.asm.Push()
	}

	for i := len(g.ast.Children(argList)) - 1; i >= 0; i-- {
		g.asm.Pop(i)
	}

	g.asm.Call(g.ast.NodeString(name))
}
