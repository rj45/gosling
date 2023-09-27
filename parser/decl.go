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

		if p.tok.Kind() == token.Semicolon {
			p.next()
		}

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

// funcDecl = "func" ident "(" fieldList? ")" ident? block
func (p *Parser) funcDecl() ast.NodeID {
	tok := p.expect(token.Func)
	name := p.name()

	p.expect(token.LParen)
	params := p.fieldList(token.Comma, token.RParen)
	p.expect(token.RParen)

	var ret ast.NodeID
	if p.tok.Kind() == token.Ident {
		ret = p.name()
	}

	body := p.block()

	return p.ast.AddNode(ast.FuncDecl, tok, name, params, ret, body)
}

// fieldList = (field (sep field)*)?
func (p *Parser) fieldList(sep token.Kind, end token.Kind) ast.NodeID {
	tok := p.tok
	if p.tok.Kind() == end {
		return p.ast.AddNode(ast.FieldList, tok)
	}

	nodes := []ast.NodeID{p.field()}
	for p.tok.Kind() == token.Comma {
		p.next()
		nodes = append(nodes, p.field())
	}
	return p.ast.AddNode(ast.FieldList, tok, nodes...)
}

// field = ident ident
func (p *Parser) field() ast.NodeID {
	tok := p.tok
	name := p.name()
	if p.tok.Kind() != token.Ident {
		p.error("expected type")
		return ast.InvalidNode
	}
	typ := p.name()
	return p.ast.AddNode(ast.Field, tok, name, typ)
}
