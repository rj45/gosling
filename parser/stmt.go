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

// stmt = exprStmt ";"
func (p *Parser) stmt() ast.NodeID {
	node := p.exprStmt()
	if p.tok.Kind() != token.EOF {
		p.expect(token.Semicolon)
	}
	return node
}

// exprStmt = expr
func (p *Parser) exprStmt() ast.NodeID {
	return p.ast.AddNode(ast.ExprStmt, p.tok, p.expr())
}
