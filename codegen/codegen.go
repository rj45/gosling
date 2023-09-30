package codegen

import (
	"github.com/rj45/gosling/ast"
)

type Assembly interface {
	WordSize() int

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
	JumpIfFalse(string, int)
	Jump(string, int)
	Label(string, int)
}

type CodeGen struct {
	ast    *ast.AST
	symtab *ast.SymTab
	asm    Assembly
	label  int
}

func New(ast *ast.AST, symtab *ast.SymTab, asm Assembly) *CodeGen {
	return &CodeGen{
		ast:    ast,
		symtab: symtab,
		asm:    asm,
	}
}

func (g *CodeGen) Generate() {
	g.genDeclList(g.ast.Root())
}
