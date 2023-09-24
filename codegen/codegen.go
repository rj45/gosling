package codegen

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

type Assembly interface {
	WordSize() int

	Prologue(string, int)
	Epilogue()

	Push()
	Pop()
	LoadLocal(int)
	Load()
	Store()

	LoadInt(string)
	LocalAddr(int)

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
	JumpIfFalse(string, int)
	Jump(string, int)
	Label(string, int)
}

type CodeGen struct {
	ast   *ast.AST
	asm   Assembly
	label int
}

func New(ast *ast.AST, asm Assembly) *CodeGen {
	return &CodeGen{
		ast: ast,
		asm: asm,
	}
}

func (g *CodeGen) Generate() {

	g.genDeclList(g.ast.Root())

}

func (g *CodeGen) genDeclList(node ast.NodeID) {
	decls := g.ast.Children(node)

	// generate main func first
	for _, decl := range decls {
		if g.ast.Kind(decl) != ast.FuncDecl {
			continue
		}
		if g.ast.NodeString(g.ast.Child(decl, ast.FuncDeclName)) != "main" {
			continue
		}
		g.genDecl(decl)
	}

	// generate other funcs
	for _, decl := range decls {
		if g.ast.Kind(decl) == ast.FuncDecl && g.ast.NodeString(g.ast.Child(decl, ast.FuncDeclName)) == "main" {
			continue
		}
		g.genDecl(decl)
	}
}

func (g *CodeGen) genDecl(node ast.NodeID) {
	switch g.ast.Kind(node) {
	case ast.FuncDecl:
		g.genFuncDecl(node)
	default:
		panic("unknown decl kind")
	}
}

func (g *CodeGen) genFuncDecl(node ast.NodeID) {
	name := g.ast.Child(node, ast.FuncDeclName)

	g.asm.Prologue(g.ast.NodeString(name), g.ast.StackSize())
	body := g.ast.Child(node, ast.FuncDeclBody)

	g.genStmtList(body)

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
	g.asm.Pop()
	g.asm.Store()
}

func (g *CodeGen) genReturnStmt(node ast.NodeID) {
	for _, child := range g.ast.Children(node) {
		g.genExpr(child)
	}
	g.asm.JumpToEpilogue()
}

func (g *CodeGen) genIfExpr(node ast.NodeID) {
	cond := g.ast.Child(node, ast.IfExprCond)
	then := g.ast.Child(node, ast.IfExprThen)
	els := g.ast.Child(node, ast.IfExprElse)

	label := g.label
	g.label++

	g.genExpr(cond)
	g.asm.JumpIfFalse("else", label)
	g.genStmt(then)
	g.asm.Jump("endif", label)
	g.asm.Label("else", label)
	if els != ast.InvalidNode {
		g.genStmt(els)
	}
	g.asm.Label("endif", label)
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

func (g *CodeGen) localOffset(node ast.NodeID) int {
	switch g.ast.Kind(node) {
	case ast.Name:
		sym := g.ast.SymbolOf(node)
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
		g.asm.Pop()

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
		g.genStmtList(node)
	case ast.Literal:
		g.asm.LoadInt(g.ast.NodeString(node))
	case ast.Name:
		sym := g.ast.SymbolOf(node)
		if sym.Const != nil {
			g.genConst(sym.Const)
			return
		}
		g.asm.LoadLocal(g.localOffset(node))
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
