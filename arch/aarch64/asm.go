package aarch64

import "fmt"

type Assembly struct {
	depth int
}

func (g *Assembly) Prologue() {
	fmt.Println(".text")
	fmt.Println(".global _main")
	fmt.Println(".align 2")
	fmt.Println("_main:")
}

func (g *Assembly) Push() {
	g.depth++
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	fmt.Println("  str x0, [sp, #-16]!")
}

func (g *Assembly) Pop() {
	g.depth--
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	fmt.Println("  ldr x1, [sp], #16")
}

func (g *Assembly) LoadInt(lit string) {
	fmt.Println("  mov x0, #" + lit)
}

func (g *Assembly) Add() {
	fmt.Println("  add x0, x1, x0")
}

func (g *Assembly) Sub() {
	fmt.Println("  sub x0, x1, x0")
}

func (g *Assembly) Mul() {
	fmt.Println("  mul x0, x1, x0")
}

func (g *Assembly) Div() {
	fmt.Println("  sdiv x0, x1, x0")
}

func (g *Assembly) Neg() {
	fmt.Println("  neg x0, x0")
}

func (g *Assembly) Epilogue() {
	fmt.Println("  mov x16, #1") // syscall number for exit()
	fmt.Println("  svc #0")      // syscall
	fmt.Println("  ret")

	if g.depth != 0 {
		panic("unbalanced stack")
	}
}
