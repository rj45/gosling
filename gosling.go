package main

import (
	"fmt"
	"os"
)

func main() {
	// Aarch64 assembly for exit() syscall
	fmt.Println(".text")
	fmt.Println(".global _main")
	fmt.Println(".align 2")
	fmt.Println("_main:")
	fmt.Println("  mov x16, #1")            // syscall number for exit()
	fmt.Println("  mov x0, #" + os.Args[1]) // exit code
	fmt.Println("  svc #0")                 // syscall
	fmt.Println("  ret")
}
