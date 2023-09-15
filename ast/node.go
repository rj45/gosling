package ast

import "github.com/rj45/gosling/token"

// node represents the most frequently
// required data about an AST node, stored
// bit-packed into a 64 bit integer
type node uint64

// newNode creates a new AST node
func newNode(kind Kind, token token.Token, firstChild int) node {
	if firstChild > 0xffffff {
		panic("too many children")
	}
	return node(uint64(kind) | uint64(firstChild)<<8 | uint64(token)<<32)
}

// kind of the AST Node
func (n node) kind() Kind {
	// least significant 8 bits
	return Kind(n & 0xff)
}

// firstChild returns the index of the first child node
func (n node) firstChild() int {
	return int((n >> 8) & 0xffffff)
}

// token returns the token for the node
func (n node) token() token.Token {
	return token.Token(n >> 32)
}
