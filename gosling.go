package main

import (
	"fmt"
	"os"

	"github.com/rj45/gosling/parser"
)

func main() {
	gen := codegen{}

	gen.GenPrologue()

	src := []byte(os.Args[1])
	parser := parser.New(src, gen)
	parser.Parse()

	gen.GenEpilogue()
}

type codegen struct{}

func (g codegen) GenPrologue() {
	// Aarch64 assembly
	fmt.Println(".text")
	fmt.Println(".global _main")
	fmt.Println(".align 2")
	fmt.Println("_main:")
}

func (g codegen) GenLoadInt(lit string) {
	fmt.Println("  mov x0, #" + lit)
}

func (g codegen) GenAddInt(lit string) {
	fmt.Println("  add x0, x0, #" + lit)
}

func (g codegen) GenSubInt(lit string) {
	fmt.Println("  sub x0, x0, #" + lit)
}

func (g codegen) GenEpilogue() {
	fmt.Println("  mov x16, #1") // syscall number for exit()
	fmt.Println("  svc #0")      // syscall
	fmt.Println("  ret")
}
