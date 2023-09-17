package aarch64

import "fmt"

type Assembly struct {
	depth int
}

func align(n int, align int) int {
	return (n + align - 1) / align * align
}

func (g *Assembly) Prologue(stacksize int) {
	fmt.Println(".text")
	fmt.Println(".global _main")
	fmt.Println(".align 2")
	fmt.Println("_main:")
	fmt.Println("  stp x29, x30, [sp, #-16]!")
	fmt.Println("  mov x29, sp")
	fmt.Println("  sub sp, sp, #" + fmt.Sprintf("%d", align(stacksize*g.WordSize(), 16)))
}

func (g *Assembly) WordSize() int {
	return 8
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

func (g *Assembly) LoadLocal(offset int) {
	fmt.Printf("  ldr x0, [x29, #%d]\n", offset)
}

func (g *Assembly) StoreLocal(offset int) {
	fmt.Printf("  str x0, [x29, #%d]\n", offset)
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

func (g *Assembly) Eq() {
	fmt.Println("  cmp x1, x0")
	fmt.Println("  cset x0, eq")
}

func (g *Assembly) Ne() {
	fmt.Println("  cmp x1, x0")
	fmt.Println("  cset x0, ne")
}

func (g *Assembly) Lt() {
	fmt.Println("  cmp x1, x0")
	fmt.Println("  cset x0, lt")
}

func (g *Assembly) Le() {
	fmt.Println("  cmp x1, x0")
	fmt.Println("  cset x0, le")
}

func (g *Assembly) Gt() {
	fmt.Println("  cmp x1, x0")
	fmt.Println("  cset x0, gt")
}

func (g *Assembly) Ge() {
	fmt.Println("  cmp x1, x0")
	fmt.Println("  cset x0, ge")
}

func (g *Assembly) JumpToEpilogue() {
	fmt.Println("  b .L.epilogue")
}

func (g *Assembly) Epilogue() {
	fmt.Println(".L.epilogue:")
	fmt.Println("  mov sp, x29")
	fmt.Println("  ldp x29, x30, [sp], #16")
	fmt.Println("  mov x16, #1") // syscall number for exit()
	fmt.Println("  svc #0")      // syscall
	fmt.Println("  ret")

	if g.depth != 0 {
		panic("unbalanced stack")
	}
}
