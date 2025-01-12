package ir

import (
	"fmt"

	"github.com/rj45/gosling/assert"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// NodeID is the index of the node in the graph
type NodeID uint32

type NodeScope uint8

const (
	InvalidNodeScope NodeScope = iota
	LocalScope
	PackageScope
	GlobalScope
)

const NodeIDIndexMask = 0x3fffffff

func NewNodeID(kind NodeScope, index int) NodeID {
	assert.True(index < NodeIDIndexMask)
	assert.True(index >= 0)
	return (NodeID(kind) << 30) | NodeID(index&NodeIDIndexMask)
}

func (nid NodeID) Scope() NodeScope {
	return NodeScope(nid >> 30)
}

func (nid NodeID) Index() int {
	return int(nid & NodeIDIndexMask)
}

const InvalidNode NodeID = 0

// Node is a fat pointer representing a node in the graph.
// Nodes are either control flow or data flow or both.
// Traditional basic blocks are represented as one or more
// nodes. In general a node has a list of inputs, an Op,
// a Type, and a list of Uses.
type Node struct {
	g  *Graph
	ID NodeID
}

func (n Node) Graph() *Graph {
	return n.g
}

func (n Node) Index() int {
	return n.ID.Index()
}

func (n Node) Scope() NodeScope {
	return n.ID.Scope()
}

func (n Node) IsValid() bool {
	return n.Scope() != InvalidNodeScope
}

func (n Node) Op() Op {
	return n.g.op(n.ID)
}

func (n Node) Flags() uint8 {
	return n.g.flags(n.ID)
}

func (n Node) Type() types.Type {
	return n.g.typ(n.ID)
}

func (n Node) Token() token.Token {
	return n.g.token(n.ID)
}

func (n Node) NumInputs() int {
	return n.g.numInputs(n.ID)
}

func (n Node) Input(i int) Node {
	return n.g.input(n.ID, i)
}

func (n Node) InputsIter() func(yield func(int, Node) bool) {
	return n.g.inputsIter(n.ID)
}

func (n Node) UseHead() Use {
	return n.g.useHead(n.ID)
}

func (n Node) UseIter() func(yield func(Use) bool) {
	return n.g.useIter(n.ID)
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

func (n Node) Name() string {
	name := n.Graph().names[n.ID.Index()]
	namestr := name.name
	if n.Graph().lastNameNum[name.name] > 1 {
		namestr += fmt.Sprintf("%d", name.num)
	}
	return namestr
}
