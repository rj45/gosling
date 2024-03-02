package ir

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/rj45/gosling/assert"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// Graph contains a Sea of Nodes. It is the top-level container for the IR.
// Both control flow and data flow are represented as nodes in the graph.
// The two flows sometimes intersect, such as at Phi nodes, or Return nodes.
type Graph struct {
	scope  NodeScope
	scopes [4]*Graph

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

	start NodeID
	end   NodeID

	// global value numbering -- essentially a cache for common sub-expression elimination
	gvn map[uint32]NodeID
}

// Init initializes the graph to a clean state, ensuring that invalid values are reserved.
func (g *Graph) Init(scope NodeScope, local *Graph, pkg *Graph, global *Graph) {
	*g = Graph{
		scope:    scope,
		scopes:   [4]*Graph{nil, local, pkg, global},
		nodes:    []node{{}},                    // node 0 is invalid
		types:    []types.Type{types.None},      // type 0 is invalid
		useHeads: []inputID{invalidInput},       // use head 0 is invalid
		tokens:   []token.Token{token.Token(0)}, // token 0 is invalid
		inputs:   []NodeID{InvalidNode},         // input 0 is invalid
		uses:     []use{{}},                     // use 0 is invalid
		ops:      []Op{OpInvalid},               // op 0 is invalid
		opNames:  []string{"invalid"},           // op 0 is invalid
		names:    []name{{"invalid", 0}},        // node 0 is invalid
	}
}

func (g *Graph) NumNodes() int {
	return len(g.nodes) - 1
}

func (g *Graph) Node(id NodeID) Node {
	return Node{ID: id, g: g.scopes[id.Scope()]}
}

func (g *Graph) nodeAt(idx int) Node {
	return g.Node(NewNodeID(g.scope, idx))
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

				if g.inputs[g.nodes[id.Index()].firstInput+inputID(i)] != input {
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

	id := g.nodeAt(len(g.nodes)).ID
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
			nextUse: g.useHeads[input.Index()],
		})
		g.useHeads[input.Index()] = inputID(len(g.uses) - 1)
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
		inputIDs = append(inputIDs, input.ID)
	}

	id := g.NewNodeWithIDs(op, typ, token, inputIDs...)
	return g.Node(id)
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
	return g.ops[g.nodes[id.Index()].op]
}

func (g *Graph) flags(id NodeID) uint8 {
	return g.nodes[id.Index()].flags
}

func (g *Graph) typ(id NodeID) types.Type {
	return g.types[id.Index()]
}

func (g *Graph) token(id NodeID) token.Token {
	return g.tokens[id.Index()]
}

func (g *Graph) numInputs(id NodeID) int {
	return int(g.nodes[id.Index()].numInputs)
}

func (g *Graph) input(id NodeID, i int) Node {
	assert.True(i >= 0)
	assert.True(i < int(g.nodes[id.Index()].numInputs))
	return g.Node(g.inputs[g.nodes[id.Index()].firstInput+inputID(i)])
}

func (g *Graph) inputsIter(id NodeID) func(yield func(int, Node) bool) {
	return func(yield func(int, Node) bool) {
		firstInput := g.nodes[id.Index()].firstInput
		numInputs := g.nodes[id.Index()].numInputs
		inputs := g.inputs
		for i := 0; i < int(numInputs); i++ {
			if !yield(i, g.Node(inputs[firstInput+inputID(i)])) {
				return
			}
		}
	}
}

func (g *Graph) useHead(id NodeID) Use {
	return Use{g, g.useHeads[id.Index()]}
}

func (g *Graph) useIter(id NodeID) func(yield func(Use) bool) {
	return func(yield func(Use) bool) {
		for i := g.useHeads[id.Index()]; i != invalidInput; i = g.uses[i].nextUse {
			if !yield(Use{g, i}) {
				return
			}
		}
	}
}

func (g *Graph) reassignName(id NodeID, namestr string) {
	oldName := g.names[id.Index()]
	if oldName.num == g.lastNameNum[oldName.name] {
		g.lastNameNum[oldName.name]--
	}
	g.lastNameNum[namestr]++
	g.names[id.Index()] = name{namestr, g.lastNameNum[namestr]}
}

func (g *Graph) Iter() func(yield func(Node) bool) {
	return func(yield func(Node) bool) {
		for i := range g.nodes {
			if i == 0 {
				continue
			}
			if !yield(g.nodeAt(i)) {
				return
			}
		}
	}
}

// DCEIter iterates over all nodes in the graph, doing Dead Code Elimination
func (g *Graph) DCEIter() func(yield func(Node) bool) {
	return func(yield func(Node) bool) {
		// First find all alive nodes, starting from the end node
		alive := make([]bool, len(g.nodes))
		g.visit(g.end, alive, func(node Node) bool {
			return true
		})

		// Then iterate over all nodes, yielding those that are alive
		for i := range g.nodes {
			if i == 0 {
				continue
			}
			if alive[i] {
				if !yield(g.nodeAt(i)) {
					return
				}
			}
		}
	}
}

func (g *Graph) visit(id NodeID, visited []bool, yield func(Node) bool) bool {
	if visited[id.Index()] {
		return true
	}
	visited[id.Index()] = true
	for i := 0; i < g.numInputs(id); i++ {
		if !g.visit(g.input(id, i).ID, visited, yield) {
			return false
		}
	}
	return yield(g.Node(id))
}

// BottomDFSIter iterates over the graph in depth-first order, starting from the end node.
func (g *Graph) BottomDFSIter() func(yield func(Node) bool) {
	return func(yield func(Node) bool) {
		visited := make([]bool, len(g.nodes))
		g.visit(g.end, visited, yield)
	}
}
