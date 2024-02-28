package ir

import (
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// NodeID is the index of the node in the graph
type NodeID uint32

const InvalidNode NodeID = 0

type Node struct {
	g  *Graph
	id NodeID
}

func (n Node) ID() NodeID {
	return n.id
}

func (n Node) Graph() *Graph {
	return n.g
}

func (n Node) Op() Op {
	return n.g.op(n.id)
}

func (n Node) Flags() uint8 {
	return n.g.flags(n.id)
}

func (n Node) Type() types.Type {
	return n.g.typ(n.id)
}

func (n Node) Token() token.Token {
	return n.g.token(n.id)
}

func (n Node) NumInputs() int {
	return n.g.numInputs(n.id)
}

func (n Node) Input(i int) Node {
	return n.g.input(n.id, i)
}

func (n Node) UseHead() Use {
	return n.g.useHead(n.id)
}

func (n Node) IsControlFlow() bool {
	return n.Op().IsControlFlow()
}

func (n Node) IsDataFlow() bool {
	return n.Op().IsDataFlow()
}

func (n Node) IsConstant() bool {
	return n.Type().Kind() == types.ConstType
}
