package parser

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

// stmtList = stmt*
func (p *Parser) stmtList() ast.NodeID {
	nodes := []ast.NodeID{}
	tok := p.tok
	for p.tok.Kind() != token.EOF {
		nodes = append(nodes, p.stmt())
	}
	return p.ast.AddNode(ast.StmtList, tok, nodes...)
}

// stmt = simpleStmt ";"
func (p *Parser) stmt() ast.NodeID {
	node := p.simpleStmt()
	if p.tok.Kind() != token.EOF {
		p.expect(token.Semicolon)
	}
	return node
}

// simpleStmt = name "=" expr | expr
func (p *Parser) simpleStmt() ast.NodeID {
	tok := p.tok

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
