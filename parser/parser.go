package parser

import (
	"fmt"
	"os"

	"github.com/rj45/gosling/token"
)

type CodeGen interface {
	GenLoadInt(string)
	GenAddInt(string)
	GenSubInt(string)
}

type Parser struct {
	src []byte
	tok token.Token
	gen CodeGen
}

func New(src []byte, gen CodeGen) *Parser {
	return &Parser{src: src, gen: gen}
}

func (p *Parser) Parse() {
	p.next()

	tok := p.expect(token.Int)
	p.gen.GenLoadInt(tok.StringFrom(p.src))

	for p.tok.Kind() != token.EOF {
		switch p.tok.Kind() {
		case token.Add:
			p.next()
			tok = p.expect(token.Int)
			p.gen.GenAddInt(tok.StringFrom(p.src))
		case token.Sub:
			p.next()
			tok = p.expect(token.Int)
			p.gen.GenSubInt(tok.StringFrom(p.src))
		default:
			p.error("expected + or -")
		}
	}
}

func (p *Parser) next() {
	p.tok = p.tok.Next(p.src)
}

// expect errors if the token is not of the expected kind, otherwise it returns
// the next token and the token text.
func (p *Parser) expect(kind token.Kind) token.Token {
	tok := p.tok
	if p.tok.Kind() != kind {
		p.error("expected %s", kind.String())
	}
	p.next()
	return tok
}

// error prints the source code and a pointer to the token that caused the error.
func (p *Parser) error(msg string, args ...interface{}) {
	// todo: print the line number and column number instead of offset
	fmt.Fprintf(os.Stderr, "Error at offset %d:\n", p.tok.Offset())
	fmt.Fprintf(os.Stderr, "%s\n", string(p.src))
	fmt.Fprintf(os.Stderr, "%*s^  ", p.tok.Offset(), "")
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
