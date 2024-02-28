package ir

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// Graph contains a Sea of Nodes. It is the top-level container for the IR.
// Both control flow and data flow are represented as nodes in the graph.
// The two flows sometimes intersect, such as at Phi nodes, or Return nodes.
type Graph struct {
	nodes []node // indexed by NodeID

	// TODO: maybe include these in node?
	types    []types.Type // indexed by NodeID
	useHeads []inputID    // indexed by NodeID

	tokens []token.Token // indexed by NodeID

	inputs []NodeID // indexed by inputID
	uses   []use    // indexed by inputID

	ops     []Op     // indexed by opID
	opNames []string // indexed by opID

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

				if g.inputs[g.nodes[id].firstInput+inputID(i)] != input {
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

	nodeOpID := invalidOp
	for i, o := range g.ops {
		if o.Opcode() == op.Opcode() {
			nodeOpID = opID(i)
			break
		}
	}
	if nodeOpID == invalidOp {
		nodeOpID = opID(len(g.ops))
		g.ops = append(g.ops, op)
		g.opNames = append(g.opNames, strings.ToLower(op.String()))
	}

	id := NodeID(len(g.nodes))
	g.nodes = append(g.nodes, node{
		op:         nodeOpID,
		flags:      0,
		numInputs:  uint16(len(inputs)),
		firstInput: inputID(len(g.inputs)),
	})
	g.tokens = append(g.tokens, token)
	g.useHeads = append(g.useHeads, invalidInput)

	for _, input := range inputs {
		g.inputs = append(g.inputs, input)
		g.uses = append(g.uses, use{
			user:    id,
			nextUse: g.useHeads[input],
		})
		g.useHeads[input] = inputID(len(g.uses) - 1)
	}

	namestr := g.opNames[nodeOpID]
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

// opID is the index of the Op in the graph's ops array
type opID uint8

const invalidOp opID = 0

// inputID is the index of the input in the inputs array
type inputID uint32

const invalidInput inputID = 0

type name struct {
	name string
	num  int
}

type node struct {
	op    opID
	flags uint8

	numInputs  uint16
	firstInput inputID
}

type use struct {
	user    NodeID
	nextUse inputID
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
	return Node{g, g.inputs[g.nodes[id].firstInput+inputID(i)]}
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
