package parser

import "github.com/rj45/gosling/token"

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
			panic("invalid token: " + p.tok.String())
		}
	}
}

func (p *Parser) next() {
	p.tok = p.tok.Next(p.src)
}

// expect panics if the token is not of the expected kind, otherwise it returns
// the next token and the token text.
func (p *Parser) expect(kind token.Kind) token.Token {
	tok := p.tok
	if p.tok.Kind() != kind {
		panic("expected " + kind.String() + " token")
	}
	p.next()
	return tok
}
