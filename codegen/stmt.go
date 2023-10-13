package codegen

import "github.com/rj45/gosling/ast"

func (g *CodeGen) genStmtList(node ast.NodeID, last bool) {
	g.symtab.EnterScope(node)
	defer g.symtab.LeaveScope()

	children := g.ast.Children(node)
	for i, child := range children {
		g.genStmt(child, last && (i == len(children)-1))
	}
}

func (g *CodeGen) genStmt(node ast.NodeID, last bool) {
	switch g.ast.Kind(node) {
	case ast.ExprStmt:
		g.genExpr(g.ast.Child(node, ast.ExprStmtExpr))
	case ast.AssignStmt:
		g.genAssignStmt(node)
	case ast.ReturnStmt:
		g.genReturnStmt(node, last)
	case ast.IfExpr:
		g.genIfExpr(node)
	case ast.ForStmt:
		g.genForStmt(node)
	case ast.StmtList:
		g.genStmtList(node, last)
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

func (g *CodeGen) genReturnStmt(node ast.NodeID, last bool) {
	for _, child := range g.ast.Children(node) {
		g.genExpr(child)
	}
	g.asm.JumpToEpilogue()
	if last {
		return
	}

	// make sure a new block is created after the return
	g.label++
	g.asm.Label("post.return", g.label)
}

func (g *CodeGen) genForStmt(node ast.NodeID) {
	init := g.ast.Child(node, ast.ForStmtInit)
	cond := g.ast.Child(node, ast.ForStmtCond)
	post := g.ast.Child(node, ast.ForStmtPost)
	body := g.ast.Child(node, ast.ForStmtBody)

	label := g.label
	g.label++

	if init != ast.InvalidNode {
		g.genStmt(init, false)
	}
	g.asm.Label("loop", label)
	if cond != ast.InvalidNode {
		g.genExpr(g.ast.Child(cond, ast.ExprStmtExpr))
		g.asm.JumpIfFalse("endloop", label)
		g.asm.Label("loopbody", label)
	}
	g.genStmt(body, false)
	if post != ast.InvalidNode {
		g.genStmt(post, false)
	}
	g.asm.Jump("loop", label)
	g.asm.Label("endloop", label)
}
