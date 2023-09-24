package aarch64

import (
	"fmt"
	"io"
)

type Assembly struct {
	Out   io.Writer
	depth int
	fn    string
}

func align(n int, align int) int {
	return (n + align - 1) / align * align
}

func (g *Assembly) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.Out, format+"\n", args...)
}

func (g *Assembly) Prologue(fnname string, stacksize int) {
	g.fn = "_" + fnname
	g.printf(".text")
	g.printf(".global %s", g.fn)
	g.printf(".align 2")
	g.printf("%s:", g.fn)
	g.printf("  stp x29, x30, [sp, #-16]!")
	g.printf("  mov x29, sp")
	g.printf("  sub sp, sp, #%d", align(stacksize*g.WordSize(), 16))
}

func (g *Assembly) WordSize() int {
	return 8
}

func (g *Assembly) Push() {
	g.depth++
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	g.printf("  str x0, [sp, #-16]!")
}

func (g *Assembly) Pop() {
	g.depth--
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	g.printf("  ldr x1, [sp], #16")
}

func (g *Assembly) LoadLocal(offset int) {
	g.printf("  ldr x0, [x29, #%d]\n", -offset)
}

func (g *Assembly) Load() {
	g.printf("  ldr x0, [x0]")
}

func (g *Assembly) Store() {
	g.printf("  str x0, [x1]")
}

func (g *Assembly) LoadInt(lit string) {
	g.printf("  mov x0, #" + lit)
}

func (g *Assembly) LocalAddr(offset int) {
	g.printf("  add x0, x29, #%d\n", -offset)
}

func (g *Assembly) Add() {
	g.printf("  add x0, x1, x0")
}

func (g *Assembly) Sub() {
	g.printf("  sub x0, x1, x0")
}

func (g *Assembly) Mul() {
	g.printf("  mul x0, x1, x0")
}

func (g *Assembly) Div() {
	g.printf("  sdiv x0, x1, x0")
}

func (g *Assembly) Neg() {
	g.printf("  neg x0, x0")
}

func (g *Assembly) Eq() {
	g.printf("  cmp x1, x0")
	g.printf("  cset x0, eq")
}

func (g *Assembly) Ne() {
	g.printf("  cmp x1, x0")
	g.printf("  cset x0, ne")
}

func (g *Assembly) Lt() {
	g.printf("  cmp x1, x0")
	g.printf("  cset x0, lt")
}

func (g *Assembly) Le() {
	g.printf("  cmp x1, x0")
	g.printf("  cset x0, le")
}

func (g *Assembly) Gt() {
	g.printf("  cmp x1, x0")
	g.printf("  cset x0, gt")
}

func (g *Assembly) Ge() {
	g.printf("  cmp x1, x0")
	g.printf("  cset x0, ge")
}

func (g *Assembly) JumpToEpilogue() {
	g.printf("  b .L.epilogue%s", g.fn)
}

func (g *Assembly) JumpIfFalse(label string, count int) {
	g.printf("  cmp x0, #0")
	g.printf("  b.eq .L.%s.%d\n", label, count)
}

func (g *Assembly) Jump(label string, count int) {
	g.printf("  b .L.%s.%d\n", label, count)
}

func (g *Assembly) Label(label string, count int) {
	g.printf(".L.%s.%d:\n", label, count)
}

func (g *Assembly) Epilogue() {
	g.printf(".L.epilogue%s:", g.fn)
	g.printf("  mov sp, x29")
	g.printf("  ldp x29, x30, [sp], #16")
	g.printf("  mov x16, #1") // syscall number for exit()
	g.printf("  svc #0")      // syscall
	g.printf("  ret")

	if g.depth != 0 {
		panic("unbalanced stack")
	}
}
