package ir

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// OpID is the index of the Op in the graph's Ops array
type OpID uint8

const InvalidOp OpID = 0

// NodeID is the index of the node in the graph
type NodeID uint32

const InvalidNode NodeID = 0

// InputID is the index of the input in the inputs array
type InputID uint32

const InvalidInput InputID = 0

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

type Use struct {
	g  *Graph
	id InputID
}

func (u Use) User() Node {
	return Node{u.g, u.g.uses[u.id].user}
}

func (u Use) Next() Use {
	return Use{u.g, u.g.uses[u.id].nextUse}
}

func (u Use) Def() Node {
	return Node{u.g, u.g.inputs[u.id]}
}

func (u Use) Head() Use {
	return Use{u.g, u.g.useHeads[u.g.inputs[u.id]]}
}

type name struct {
	name string
	num  int
}

type Graph struct {
	nodes []node // indexed by NodeID

	// TODO: maybe include these in node?
	types    []types.Type // indexed by NodeID
	useHeads []InputID    // indexed by NodeID

	tokens []token.Token // indexed by NodeID

	inputs []NodeID // indexed by InputID
	uses   []use    // indexed by InputID

	ops     []Op     // indexed by OpID
	opNames []string // indexed by OpID

	names       []name // indexed by NodeID
	lastNameNum map[string]int

	// global value numbering -- essentially a cache for common sub-expression elimination
	gvn map[uint32]NodeID
}

func (g *Graph) NumNodes() int {
	return len(g.nodes) - 1
}

func (g *Graph) Node(id NodeID) Node {
	return Node{g, id}
}

func (g *Graph) NewNodeWithIDs(op Op, typ types.Type, token token.Token, inputs ...NodeID) NodeID {
	// TODO: this needs to be optimized

	if reflect.ValueOf(op).Type().Kind() == reflect.Ptr {
		panic("op is a pointer")
	}

	code := hashCode(uint32(op.Opcode()), inputs...)
	if id, ok := g.gvn[code]; ok {
		if g.op(id) == op && g.typ(id) == typ {
			found := true
			for i, input := range inputs {

				if g.inputs[g.nodes[id].firstInput+InputID(i)] != input {
					found = false
					break
				}
			}
			if found {
				return id
			}
		}
		// hash collision
		fmt.Println("hash collision:", id, op, typ, inputs)
	}

	g.types = append(g.types, typ)

	opID := InvalidOp
	for i, o := range g.ops {
		if o.Opcode() == op.Opcode() {
			opID = OpID(i)
			break
		}
	}
	if opID == InvalidOp {
		opID = OpID(len(g.ops))
		g.ops = append(g.ops, op)
		g.opNames = append(g.opNames, strings.ToLower(op.String()))
	}

	id := NodeID(len(g.nodes))
	g.nodes = append(g.nodes, node{
		op:         opID,
		flags:      0,
		numInputs:  uint16(len(inputs)),
		firstInput: InputID(len(g.inputs)),
	})
	g.tokens = append(g.tokens, token)
	g.useHeads = append(g.useHeads, InvalidInput)

	for _, input := range inputs {
		g.inputs = append(g.inputs, input)
		g.uses = append(g.uses, use{
			user:    id,
			nextUse: g.useHeads[input],
		})
		g.useHeads[input] = InputID(len(g.uses) - 1)
	}

	namestr := g.opNames[opID]
	if g.lastNameNum == nil {
		g.lastNameNum = make(map[string]int)
	}
	g.lastNameNum[namestr]++
	g.names = append(g.names, name{namestr, g.lastNameNum[namestr]})

	return id
}

func (g *Graph) NewNode(op Op, typ types.Type, token token.Token, inputs ...Node) Node {
	inputIDs := make([]NodeID, len(inputs))

	for _, input := range inputs {
		inputIDs = append(inputIDs, input.ID())
	}

	id := g.NewNodeWithIDs(op, typ, token, inputIDs...)
	return Node{g, id}
}

type node struct {
	op    OpID
	flags uint8

	numInputs  uint16
	firstInput InputID
}

type use struct {
	user    NodeID
	nextUse InputID
}

func (g *Graph) op(id NodeID) Op {
	return g.ops[g.nodes[id].op]
}

func (g *Graph) flags(id NodeID) uint8 {
	return g.nodes[id].flags
}

func (g *Graph) typ(id NodeID) types.Type {
	return g.types[id]
}

func (g *Graph) token(id NodeID) token.Token {
	return g.tokens[id]
}

func (g *Graph) numInputs(id NodeID) int {
	return int(g.nodes[id].numInputs)
}

func (g *Graph) input(id NodeID, i int) Node {
	if i < 0 || i >= int(g.nodes[id].numInputs) {
		panic("input index out of range")
	}
	return Node{g, g.inputs[g.nodes[id].firstInput+InputID(i)]}
}

func (g *Graph) useHead(id NodeID) Use {
	return Use{g, g.useHeads[id]}
}

func (g *Graph) reassignName(id NodeID, namestr string) {
	oldName := g.names[id]
	if oldName.num == g.lastNameNum[oldName.name] {
		g.lastNameNum[oldName.name]--
	}
	g.lastNameNum[namestr]++
	g.names[id] = name{namestr, g.lastNameNum[namestr]}
}

type Program struct {
	Packages []Package

	Types *types.Universe
}

type Package struct {
	*Program

	Funcs []Function
}

type Function struct {
	Graph

	*Package

	Name string
	Type types.Type

	Start NodeID
	End   NodeID
}
