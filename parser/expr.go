package parser

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

// expr = mul ("+" mul | "-" mul)*
func (p *Parser) expr() ast.NodeID {
	node := p.mul()

	for {
		switch p.tok.Kind() {
		case token.Add, token.Sub:
			node = p.ast.AddNode(ast.BinaryExpr, p.next(), node, p.mul())
		default:
			return node
		}
	}
}

// mul = unary ("*" unary | "/" unary)*
func (p *Parser) mul() ast.NodeID {
	node := p.unary()
	for {
		switch p.tok.Kind() {
		case token.Mul, token.Div:
			node = p.ast.AddNode(ast.BinaryExpr, p.next(), node, p.unary())
		default:
			return node
		}

	}
}

// unary = ("+" | "-") unary
//
//	| primary
func (p *Parser) unary() ast.NodeID {
	switch p.tok.Kind() {
	case token.Add:
		p.next()
		return p.unary()
	case token.Sub:
		return p.ast.AddNode(ast.UnaryExpr, p.next(), p.unary())
	default:
		return p.primary()
	}
}

// primary = "(" expr ")" | number
func (p *Parser) primary() ast.NodeID {
	switch p.tok.Kind() {
	case token.LParen:
		p.next()
		expr := p.expr()
		p.expect(token.RParen)
		return expr
	case token.Int:
		return p.node(ast.Literal, token.Int)
	default:
		p.error("expected int")
		return ast.InvalidNode
	}
}
