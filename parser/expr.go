package parser

import (
	"github.com/rj45/gosling/ast"
	"github.com/rj45/gosling/token"
)

// expr = equality
func (p *Parser) expr() ast.NodeID {
	return p.equality()
}

// equality = comparison ("==" comparison | "!=" comparison)*
func (p *Parser) equality() ast.NodeID {
	node := p.comparison()

	for {
		switch p.tok.Kind() {
		case token.Eq, token.Ne:
			node = p.ast.AddNode(ast.BinaryExpr, p.next(), node, p.comparison())
		default:
			return node
		}
	}
}

// comparison = add ("<" add | "<=" add | ">" add | ">=" add)*
func (p *Parser) comparison() ast.NodeID {
	node := p.add()

	for {
		switch p.tok.Kind() {
		case token.Lt, token.Le, token.Gt, token.Ge:
			node = p.ast.AddNode(ast.BinaryExpr, p.next(), node, p.add())
		default:
			return node
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func (p *Parser) add() ast.NodeID {
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
		case token.Star, token.Div:
			node = p.ast.AddNode(ast.BinaryExpr, p.next(), node, p.unary())
		default:
			return node
		}
	}
}

// unary = ("+" | "-" | "*" | "&") unary
//
//	| primary
func (p *Parser) unary() ast.NodeID {
	switch p.tok.Kind() {
	case token.Add:
		p.next()
		return p.unary()
	case token.Sub:
		return p.ast.AddNode(ast.UnaryExpr, p.next(), p.unary())
	case token.Star:
		return p.ast.AddNode(ast.DerefExpr, p.next(), p.unary())
	case token.And:
		return p.ast.AddNode(ast.AddrExpr, p.next(), p.unary())
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
	case token.LBrace:
		return p.block()
	case token.If:
		return p.ifExpr()
	case token.Int:
		return p.node(ast.Literal, token.Int)
	case token.Ident:
		node := p.node(ast.Name, token.Ident)
		p.ast.SymbolOf(node) // generate address
		if p.tok.Kind() == token.LParen {
			return p.ast.AddNode(ast.CallExpr, p.tok, node, p.argList())
		}
		return node
	default:
		p.error("expected expression")
		return ast.InvalidNode
	}
}

// argList = "(" (expr ("," expr)*)? ")"
func (p *Parser) argList() ast.NodeID {
	p.expect(token.LParen)
	if p.tok.Kind() == token.RParen {
		p.next()
		return p.ast.AddNode(ast.ExprList, p.tok)
	}
	args := p.exprList()
	p.expect(token.RParen)
	return args
}

// exprList = expr ("," expr)*
func (p *Parser) exprList() ast.NodeID {
	tok := p.tok
	nodes := []ast.NodeID{p.expr()}
	for p.tok.Kind() == token.Comma {
		p.next()
		nodes = append(nodes, p.expr())
	}
	return p.ast.AddNode(ast.ExprList, tok, nodes...)
}

// block = "{" stmtList "}"
func (p *Parser) block() ast.NodeID {
	p.expect(token.LBrace)
	stmts := p.stmtList()
	p.expect(token.RBrace)
	return stmts
}

// ifExpr = "if" expr blockStmt ("else" blockStmt)?
func (p *Parser) ifExpr() ast.NodeID {
	tok := p.expect(token.If)
	cond := p.expr()
	then := p.block()
	if p.tok.Kind() != token.Else {
		return p.ast.AddNode(ast.IfExpr, tok, cond, then, ast.InvalidNode)
	}
	p.expect(token.Else)
	return p.ast.AddNode(ast.IfExpr, tok, cond, then, p.block())
}
