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
		nodes = append(nodes, p.stmt())
		if len(p.errs) > 0 {
			// todo: implement error recovery
			break
		}
	}
	return p.ast.AddNode(ast.StmtList, tok, nodes...)
}

// stmt = returnStmt | ifStmt | forStmt | blockStmt | simpleStmt
func (p *Parser) stmt() ast.NodeID {
	switch p.tok.Kind() {
	case token.Return:
		stmt := p.returnStmt()
		if p.tok.Kind() != token.RBrace {
			p.expect(token.Semicolon)
		}
		return stmt
	case token.If:
		return p.ifExpr()
	case token.For:
		return p.forStmt()
	case token.LBrace:
		return p.block()
	default:
		stmt := p.simpleStmt()
		if p.tok.Kind() != token.RBrace {
			p.expect(token.Semicolon)
		}
		return stmt
	}
}

// forStmt = "for" simpleStmt ";" expr ";" simpleStmt blockStmt
func (p *Parser) forStmt() ast.NodeID {
	tok := p.expect(token.For)
	var init, cond, post ast.NodeID
	if p.tok.Kind() != token.LBrace {
		init = p.simpleStmt()
		if p.tok.Kind() != token.LBrace {
			p.expect(token.Semicolon)
		}
	}

	if p.tok.Kind() != token.LBrace {
		cond = p.simpleStmt()
		if p.tok.Kind() != token.LBrace {
			p.expect(token.Semicolon)
		}
	}

	if p.tok.Kind() != token.LBrace {
		post = p.simpleStmt()
	}

	body := p.block()

	if init != ast.InvalidNode && cond == ast.InvalidNode && post == ast.InvalidNode {
		cond = init
		init = ast.InvalidNode
	}

	if cond != ast.InvalidNode && p.ast.Kind(cond) != ast.ExprStmt {
		p.errorAt(p.ast.Token(cond), "expected for condition to be expression statement")
	}

	return p.ast.AddNode(ast.ForStmt, tok, init, cond, post, body)
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
		return p.ast.AddNode(ast.EmptyStmt, p.tok)
	}

	lhs := p.expr()

	switch p.tok.Kind() {
	case token.Assign, token.Define:
		node := lhs
		for p.ast.Kind(node) == ast.DerefExpr {
			node = p.ast.Child(node, ast.DerefExprExpr)
		}

		if p.ast.Kind(node) != ast.Name {
			p.errorAt(p.ast.Token(lhs), "expected name or deref on the left side of the assignment")
		}
		return p.ast.AddNode(ast.AssignStmt, p.next(), lhs, p.expr())
	}

	return p.ast.AddNode(ast.ExprStmt, tok, lhs)
}
