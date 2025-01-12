package codegen

import "github.com/rj45/gosling/ast"

func (g *CodeGen) genStmtList(node ast.NodeID) {
	g.symtab.EnterScope(node)
	defer g.symtab.LeaveScope()

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
	case ast.IfExpr:
		g.genIfExpr(node)
	case ast.ForStmt:
		g.genForStmt(node)
	case ast.StmtList:
		g.genStmtList(node)
	case ast.EmptyStmt:
		// do nothing
	default:
		panic("unknown stmt kind")
	}
}

func (g *CodeGen) genAssignStmt(node ast.NodeID) {
	g.genAddr(g.ast.Child(node, ast.AssignStmtLHS))
	g.asm.Push()
	g.genExpr(g.ast.Child(node, ast.AssignStmtRHS))
	g.asm.Pop(1)
	g.asm.Store()
}

func (g *CodeGen) genReturnStmt(node ast.NodeID) {
	for _, child := range g.ast.Children(node) {
		g.genExpr(child)
	}
	g.asm.JumpToEpilogue()
}

func (g *CodeGen) genForStmt(node ast.NodeID) {
	init := g.ast.Child(node, ast.ForStmtInit)
	cond := g.ast.Child(node, ast.ForStmtCond)
	post := g.ast.Child(node, ast.ForStmtPost)
	body := g.ast.Child(node, ast.ForStmtBody)

	label := g.label
	g.label++

	if init != ast.InvalidNode {
		g.genStmt(init)
	}
	g.asm.Label("loop", label)
	if cond != ast.InvalidNode {
		g.genExpr(g.ast.Child(cond, ast.ExprStmtExpr))
		g.asm.JumpIfFalse("endloop", label)
	}
	g.genStmt(body)
	if post != ast.InvalidNode {
		g.genStmt(post)
	}
	g.asm.Jump("loop", label)
	g.asm.Label("endloop", label)
}
