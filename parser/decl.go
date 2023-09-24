package parser

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

// declList = decl*
func (p *Parser) declList() ast.NodeID {
	tok := p.tok

	var decls []ast.NodeID
	for p.tok.Kind() != token.EOF {
		before := p.tok
		decls = append(decls, p.decl())

		if before == p.tok {
			// do not loop forever if we aren't making progress
			p.error("expected declaration")
			return ast.InvalidNode
		}
	}

	return p.ast.AddNode(ast.DeclList, tok, decls...)
}

// decl = funcDecl
func (p *Parser) decl() ast.NodeID {
	switch p.tok.Kind() {
	case token.Func:
		return p.funcDecl()
	default:
		p.error("expected declaration")
		return ast.InvalidNode
	}
}

// funcDecl = "func" ident "(" paramList? ")" ident? block
func (p *Parser) funcDecl() ast.NodeID {
	tok := p.expect(token.Func)
	name := p.name()

	p.expect(token.LParen)
	// params := p.paramList()
	params := ast.InvalidNode
	p.expect(token.RParen)

	var ret ast.NodeID
	if p.tok.Kind() == token.Ident {
		ret = p.name()
	}

	body := p.block()

	return p.ast.AddNode(ast.FuncDecl, tok, name, params, ret, body)
}
