package parser

import (
	"fmt"

	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

type Parser struct {
	src []byte
	tok token.Token
	ast *ast.AST

	errs []ParseErr
}

func New(src []byte) *Parser {
	return &Parser{src: src, ast: ast.New(src)}
}

func (p *Parser) Parse() (*ast.AST, error) {
	p.next()

	p.blockStmt()

	var err error
	if len(p.errs) > 0 {
		err = p.errs[0]
	}

	return p.ast, err
}

func (p *Parser) next() token.Token {
	tok := p.tok
	p.tok = p.tok.Next(p.src)

	if p.tok.Kind() == token.Illegal {
		p.errorIllegalToken()
	}

	return tok
}

func (p *Parser) node(nk ast.Kind, tk token.Kind, children ...ast.NodeID) ast.NodeID {
	return p.ast.AddNode(nk, p.expect(tk), children...)
}

// expect errors if the token is not of the expected kind, otherwise it returns
// the next token and the token text.
func (p *Parser) expect(kind token.Kind) token.Token {
	if p.tok.Kind() != kind {
		p.error("expected %s", kind)
	}

	return p.next()
}

// error prints the source code and a pointer to the token that caused the error.
func (p *Parser) error(msg string, args ...interface{}) {
	p.errorAt(p.tok, msg, args...)
}

// error prints the source code and a pointer to the token that caused the error.
func (p *Parser) errorAt(tok token.Token, msg string, args ...interface{}) {
	for i, arg := range args {
		if kind, ok := arg.(token.Kind); ok {
			if str, ok := tokenErrStrs[kind]; ok {
				args[i] = str
			}
		}
	}
	p.errs = append(p.errs, ParseErr{p, tok, fmt.Sprintf(msg, args...)})
}

func (p *Parser) errorIllegalToken() {
	tok := p.tok
	pos := p.tok.Offset()
	p.tok = p.tok.NextValidToken(p.src)
	p.errorAt(tok, "illegal token %q", p.src[pos:p.tok.Offset()])
}
