package codegen

import "github.com/rj45/gosling/ast"

func (g *CodeGen) genDeclList(node ast.NodeID) {
	g.symtab.EnterScope(node)
	defer g.symtab.LeaveScope()

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
	g.symtab.EnterScope(node)
	defer g.symtab.LeaveScope()

	name := g.ast.Child(node, ast.FuncDeclName)

	g.asm.Prologue(g.ast.NodeString(name), g.symtab.StackSize())

	paramList := g.ast.Child(node, ast.FuncDeclParams)
	for i, param := range g.ast.Children(paramList) {
		name := g.ast.Child(param, ast.FieldName)
		sym := g.symtab.Lookup(g.ast.NodeString(name))
		if i != sym.Offset {
			panic("param offset mismatch")
		}
		g.asm.StoreLocal(i)
	}

	if g.symtab.StackSize() < g.ast.NumChildren(paramList) {
		panic("local size mismatch")
	}

	body := g.ast.Child(node, ast.FuncDeclBody)

	g.genStmtList(body)

	g.asm.Epilogue()
}
