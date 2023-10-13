package aarch64

import (
	"fmt"
	"io"

	"github.com/rj45/gosling/ir"
)

var argRegs = []string{"x0", "x1", "x2", "x3", "x4", "x5", "x6", "x7"}

type Assembler struct {
	Out   io.Writer
	depth int
	fn    string
}

const WordSize = 8

func align(n int, align int) int {
	return (n + align - 1) / align * align
}

func (g *Assembler) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.Out, format+"\n", args...)
}

func (g *Assembler) regFor(reg ir.RegMask) string {
	if len(reg.Regs()) != 1 {
		panic("reg must have one register")
	}
	return argRegs[reg.Pop()]
}

func (g *Assembler) Prologue(fnname string, locals int) {
	g.fn = "_" + fnname
	g.printf(".text")
	g.printf(".global %s", g.fn)
	g.printf(".align 2")
	g.printf("%s:", g.fn)
	g.printf("  stp x29, x30, [sp, #-16]!")
	g.printf("  mov x29, sp")
	g.printf("  sub sp, sp, #%d", align(locals*WordSize, 16))
}

func (g *Assembler) Epilogue() {
	g.printf("  mov sp, x29")
	g.printf("  ldp x29, x30, [sp], #16")
}

func (g *Assembler) Push(src ir.RegMask) {
	g.depth++
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	g.printf("  str %s, [sp, #-16]!", g.regFor(src))
}

func (g *Assembler) Pop(dst ir.RegMask) {
	g.depth--
	// note: stack is 16-byte aligned, so this is wasteful, we will fix that later
	g.printf("  ldr %s, [sp], #16", g.regFor(dst))
}

func (g *Assembler) LoadLocal(dst ir.RegMask, local int) {
	g.printf("  ldr %s, [x29, #%d]", g.regFor(dst), -((local + 1) * WordSize))
}

func (g *Assembler) StoreLocal(src ir.RegMask, local int) {
	g.printf("  str %s, [x29, #%d]", g.regFor(src), -((local + 1) * WordSize))
}

func (g *Assembler) Load(dst ir.RegMask, addr ir.RegMask) {
	g.printf("  ldr %s, [%s]", g.regFor(dst), g.regFor(addr))
}

func (g *Assembler) Store(src ir.RegMask, addr ir.RegMask) {
	g.printf("  str %s, [%s]", g.regFor(src), g.regFor(addr))
}

func (g *Assembler) LoadInt(dst ir.RegMask, lit int64) {
	g.printf("  mov %s, #%d", g.regFor(dst), lit)
}

func (g *Assembler) LocalAddr(dst ir.RegMask, offset int) {
	g.printf("  add %s, x29, #%d\n", g.regFor(dst), -(offset*WordSize + 8))
}

func (g *Assembler) Add(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  add %s, %s, %s", g.regFor(dst), g.regFor(src1), g.regFor(src2))
}

func (g *Assembler) Sub(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  sub %s, %s, %s", g.regFor(dst), g.regFor(src1), g.regFor(src2))
}

func (g *Assembler) Mul(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  mul %s, %s, %s", g.regFor(dst), g.regFor(src1), g.regFor(src2))
}

func (g *Assembler) Div(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  sdiv %s, %s, %s", g.regFor(dst), g.regFor(src1), g.regFor(src2))
}

func (g *Assembler) Neg(dst ir.RegMask, src ir.RegMask) {
	g.printf("  neg %s, %s", g.regFor(dst), g.regFor(src))
}

func (g *Assembler) Eq(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  cmp %s, %s", g.regFor(src1), g.regFor(src2))
	g.printf("  cset %s, eq", g.regFor(dst))
}

func (g *Assembler) Ne(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  cmp %s, %s", g.regFor(src1), g.regFor(src2))
	g.printf("  cset %s, ne", g.regFor(dst))
}

func (g *Assembler) Lt(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  cmp %s, %s", g.regFor(src1), g.regFor(src2))
	g.printf("  cset %s, lt", g.regFor(dst))
}

func (g *Assembler) Le(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  cmp %s, %s", g.regFor(src1), g.regFor(src2))
	g.printf("  cset %s, le", g.regFor(dst))
}

func (g *Assembler) Gt(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  cmp %s, %s", g.regFor(src1), g.regFor(src2))
	g.printf("  cset %s, gt", g.regFor(dst))
}

func (g *Assembler) Ge(dst ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	g.printf("  cmp %s, %s", g.regFor(src1), g.regFor(src2))
	g.printf("  cset %s, ge", g.regFor(dst))
}

func (g *Assembler) Call(fnname string) {
	g.printf("  bl _%s", fnname)
}

func (g *Assembler) If(reg ir.RegMask, then string, els string) {
	g.printf("  cmp %s, #0", g.regFor(reg))
	g.printf("  b.eq .L.%s", els)
	g.printf("  b .L.%s", then)
}

func (g *Assembler) Jump(label string) {
	g.printf("  b .L.%s", label)
}

func (g *Assembler) Label(label string) {
	g.printf(".L.%s:", label)
}

func (g *Assembler) Return() {
	if g.fn == "_main" {
		g.printf("  mov x16, #1") // syscall number for exit()
		g.printf("  svc #0")      // syscall
	}
	g.printf("  ret")

	if g.depth != 0 {
		panic("unbalanced stack")
	}
}
