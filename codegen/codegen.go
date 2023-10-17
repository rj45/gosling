package codegen

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/types"
)

type Assembly interface {
	WordSize() int

	Types(*types.Universe)

	Prologue(string, int)
	Epilogue()

	Push()
	Pop(int)
	LoadLocal(int)
	StoreLocal(int)
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

	Call(string)
	JumpToEpilogue()
	JumpIf(string, string, int)
	Jump(string, int)
	Label(string, int)

	DeclareFunction(string, types.Type)
}

type CodeGen struct {
	ast    *ast.AST
	symtab *ast.SymTab
	asm    Assembly
	types  *types.Universe
	label  int
}

func New(ast *ast.AST, symtab *ast.SymTab, types *types.Universe, asm Assembly) *CodeGen {
	return &CodeGen{
		ast:    ast,
		symtab: symtab,
		types:  types,
		asm:    asm,
	}
}

func (g *CodeGen) Generate() {
	g.asm.Types(g.types)
	g.genDeclList(g.ast.Root())
}
