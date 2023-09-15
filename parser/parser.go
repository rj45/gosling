package parser

import (
	"fmt"
	"os"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

type Parser struct {
	src []byte
	tok token.Token
	ast *ast.AST
}

func New(src []byte) *Parser {
	return &Parser{src: src, ast: ast.New(src)}
}

func (p *Parser) Parse() *ast.AST {
	p.next()

	p.expr()

	return p.ast
}

func (p *Parser) next() token.Token {
	tok := p.tok
	p.tok = p.tok.Next(p.src)
	return tok
}

func (p *Parser) node(nk ast.Kind, tk token.Kind, children ...ast.NodeID) ast.NodeID {
	return p.ast.AddNode(nk, p.expect(tk), children...)
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
