package main

import (
	"fmt"
	"os"

	"github.com/rj45/gosling/token"
)

func main() {
	// Aarch64 assembly for exit() syscall
	fmt.Println(".text")
	fmt.Println(".global _main")
	fmt.Println(".align 2")
	fmt.Println("_main:")

	src := []byte(os.Args[1])
	tok := token.Token(0)
	tok = tok.Next(src)
	var lit string

	tok, lit = expect(src, tok, token.Int)
	fmt.Println("  mov x0, #" + lit)
	for tok.Kind() != token.EOF {
		switch tok.Kind() {
		case token.Add:
			tok, lit = expect(src, tok, token.Add)
			tok, lit = expect(src, tok, token.Int)
			fmt.Println("  add x0, x0, #" + lit)
		case token.Sub:
			tok, lit = expect(src, tok, token.Sub)
			tok, lit = expect(src, tok, token.Int)
			fmt.Println("  sub x0, x0, #" + lit)
		default:
			panic("invalid token: " + tok.String())
		}
	}

	fmt.Println("  mov x0, #" + os.Args[1]) // exit code
	fmt.Println("  mov x16, #1")            // syscall number for exit()
	fmt.Println("  svc #0")                 // syscall
	fmt.Println("  ret")
}

// expect panics if the token is not of the expected kind, otherwise it returns
// the next token and the token text.
func expect(src []byte, tok token.Token, kind token.Kind) (token.Token, string) {
	if tok.Kind() != kind {
		panic("expected " + kind.String() + " token")
	}
	return tok.Next(src), tok.StringFrom(src)
}
