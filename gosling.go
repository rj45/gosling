package main

import (
	"fmt"
	"os"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/parser"
	"github.com/rj45/gosling/token"
)

func main() {
	src := []byte(os.Args[1])
	parser := parser.New(src)
	ast := parser.Parse()

	gencode(ast, &codegen{})
}

type CodeGen interface {
	Prologue()
	Epilogue()

	Push()
	Pop()

	LoadInt(string)

	Add()
	Sub()
	Mul()
	Div()
}

func gencode(ast *ast.AST, gen CodeGen) {
	gen.Prologue()
	genExpr(ast, ast.Root(), gen)
	gen.Epilogue()
}

func genExpr(a *ast.AST, node ast.NodeID, gen CodeGen) {
	switch a.Kind(node) {
	case ast.BinaryExpr:
		genExpr(a, a.Child(node, ast.BinaryExprLHS), gen)
		gen.Push()
		genExpr(a, a.Child(node, ast.BinaryExprRHS), gen)
		gen.Pop()

		switch a.Token(node).Kind() {
		case token.Add:
			gen.Add()
		case token.Sub:
			gen.Sub()
		case token.Mul:
			gen.Mul()
		case token.Div:
			gen.Div()
		}
	case ast.Literal:
		gen.LoadInt(a.NodeString(node))
	default:
		panic("unknown expr kind")
	}
}

type codegen struct {
	depth int
}

func (g *codegen) Prologue() {
	// Aarch64 assembly
	fmt.Println(".text")
	fmt.Println(".global _main")
	fmt.Println(".align 2")
	fmt.Println("_main:")
}

func (g *codegen) Push() {
	g.depth++
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	fmt.Println("  str x0, [sp, #-16]!")
}

func (g *codegen) Pop() {
	g.depth--
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	fmt.Println("  ldr x1, [sp], #16")
}

func (g *codegen) LoadInt(lit string) {
	fmt.Println("  mov x0, #" + lit)
}

func (g *codegen) Add() {
	fmt.Println("  add x0, x1, x0")
}

func (g *codegen) Sub() {
	fmt.Println("  sub x0, x1, x0")
}

func (g *codegen) Mul() {
	fmt.Println("  mul x0, x1, x0")
}

func (g *codegen) Div() {
	fmt.Println("  sdiv x0, x1, x0")
}

func (g *codegen) Epilogue() {
	fmt.Println("  mov x16, #1") // syscall number for exit()
	fmt.Println("  svc #0")      // syscall
	fmt.Println("  ret")

	if g.depth != 0 {
		panic("unbalanced stack")
	}
}
