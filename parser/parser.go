package parser

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/errors"
	"github.com/rj45/gosling/token"
)

type Parser struct {
	ast *ast.AST

	tok token.Token

	errs []error
}

func New(file *ast.File) *Parser {
	return &Parser{ast: ast.New(file)}
}

func (p *Parser) Parse() (*ast.AST, []error) {
	p.next()

	p.declList()

	return p.ast, p.errs
}

func (p *Parser) next() token.Token {
	tok := p.tok
	p.tok = p.tok.Next(p.ast.Src)

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
	p.errs = append(p.errs, errors.Newf(p.ast.Src, tok, msg, args...))
}

func (p *Parser) errorIllegalToken() {
	tok := p.tok
	pos := p.tok.Offset()
	p.tok = p.tok.NextValidToken(p.ast.Src)
	p.errorAt(tok, "illegal token %q", p.ast.Src[pos:p.tok.Offset()])
}

var tokenErrStrs = map[token.Kind]string{
	token.EOF:       "end of file",
	token.Semicolon: "newline or ';' or '}'",
	token.RBrace:    "'}'",
}
