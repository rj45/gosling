package ircodegen

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/ir"
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
	program *ir.Program
	pkg     *ir.Package
	fn      *ir.Function
	symtab  *ast.SymTab
	asm     Assembly
	label   int
	visited []bool
}

func New(ir *ir.Program, symtab *ast.SymTab, asm Assembly) *CodeGen {
	return &CodeGen{
		program: ir,
		symtab:  symtab,
		asm:     asm,
	}
}

func (g *CodeGen) Generate() {
	for pkgi := range g.program.Packages {
		pkg := &g.program.Packages[pkgi]
		for fni := range pkg.Funcs {
			fn := &pkg.Funcs[fni]
			if fn.Name == "main" {
				g.pkg = pkg
				g.fn = fn
				g.genFunction()
				break
			}
		}
	}
	for pkgi := range g.program.Packages {
		g.pkg = &g.program.Packages[pkgi]
		for fni := range g.pkg.Funcs {
			g.fn = &g.pkg.Funcs[fni]
			if g.fn.Name != "main" {
				g.genFunction()
			}
		}
	}
}
