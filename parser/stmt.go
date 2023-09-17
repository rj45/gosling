package parser

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

// stmtList = stmt*
func (p *Parser) stmtList() ast.NodeID {
	nodes := []ast.NodeID{}
	tok := p.tok
	for p.tok.Kind() != token.EOF && p.tok.Kind() != token.RBrace {
		if p.tok.Kind() == token.LBrace {
			nodes = append(nodes, p.blockStmt())
			continue
		}

		stmt := p.stmt()
		if stmt != ast.InvalidNode {
			nodes = append(nodes, stmt)
		}

		if p.tok.Kind() != token.Semicolon {
			break
		}
		p.expect(token.Semicolon)
	}
	return p.ast.AddNode(ast.StmtList, tok, nodes...)
}

// stmt = returnStmt | simpleStmt
func (p *Parser) stmt() ast.NodeID {
	switch p.tok.Kind() {
	case token.Return:
		return p.returnStmt()
	case token.LBrace:
		return p.blockStmt()
	default:
		return p.simpleStmt()
	}
}

// blockStmt = "{" stmtList "}"
func (p *Parser) blockStmt() ast.NodeID {
	p.expect(token.LBrace)
	stmts := p.stmtList()
	p.expect(token.RBrace)
	return stmts
}

// returnStmt = "return" expr?
func (p *Parser) returnStmt() ast.NodeID {
	tok := p.expect(token.Return)
	if p.tok.Kind() == token.Semicolon || p.tok.Kind() == token.EOF || p.tok.Kind() == token.RBrace {
		return p.ast.AddNode(ast.ReturnStmt, tok)
	}
	return p.ast.AddNode(ast.ReturnStmt, tok, p.expr())
}

// simpleStmt = name "=" expr | expr
func (p *Parser) simpleStmt() ast.NodeID {
	tok := p.tok

	if tok.Kind() == token.Semicolon {
		return ast.InvalidNode
	}

	lhs := p.expr()

	switch p.tok.Kind() {
	case token.Assign:
		if p.ast.Kind(lhs) != ast.Name {
			p.error("expected name on the left side of the assignment")
		}
		return p.ast.AddNode(ast.AssignStmt, p.next(), lhs, p.expr())
	}

	return p.ast.AddNode(ast.ExprStmt, tok, lhs)
}
