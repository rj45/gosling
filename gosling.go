package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	// Aarch64 assembly for exit() syscall
	fmt.Println(".text")
	fmt.Println(".global _main")
	fmt.Println(".align 2")
	fmt.Println("_main:")

	// the simplest possible parser
	prog := os.Args[1]
	tokens := strings.Fields(prog) // split by whitespace
	fmt.Println("  mov x0, #" + tokens[0])
	tokens = tokens[1:]
	for len(tokens) > 1 {
		if tokens[0] == "+" {
			fmt.Println("  add x0, x0, #" + tokens[1])
			tokens = tokens[2:]
		} else if tokens[0] == "-" {
			fmt.Println("  sub x0, x0, #" + tokens[1])
			tokens = tokens[2:]
		} else {
			panic("invalid token: " + tokens[0])
		}
	}

	fmt.Println("  mov x0, #" + os.Args[1]) // exit code
	fmt.Println("  mov x16, #1")            // syscall number for exit()
	fmt.Println("  svc #0")                 // syscall
	fmt.Println("  ret")
}
