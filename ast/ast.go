package ast

import (
	"fmt"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// NodeID is used to identify an AST node.
type NodeID uint32

// InvalidNode indicates an invalid / null node
const InvalidNode NodeID = 0

// AST is the abstract syntax tree. It is a tree of nodes,
// where each node has a kind, a token, and a list of child nodes.
// The token can be used to get the text of the node, as well as
// the position in the source code.
//
// This AST is designed to be as compact as possible in order
// to efficiently utilize the CPU cache. NodeIDs are used to
// identify nodes, and are simply an index into the node list.
// This NodeID can be used in other arrays to store additional
// information about the node.
//
// If a node has a constant number of children, it is preferable
// to identify child indices with a constant so the code is more
// readable.
type AST struct {
	File

	// node is the list of nodes in post-order traversal order
	node []node

	// child is the list of child nodes for each node indexed by
	// node.FirstChild()
	child []NodeID

	// typ is the type of each node indexed by NodeID
	typ []types.Type

	symtab   *types.SymTab
	nextAddr int
}

// New creates a new AST from the source code
func New(file *File) *AST {
	return &AST{
		File: *file,
		node: []node{
			// the first node is always an illegal node
			newNode(IllegalNode, 0, 0),
		},
		symtab: types.NewSymTab(),
	}
}

// Root returns the root node of the AST
func (a *AST) Root() NodeID {
	return NodeID(len(a.node) - 1)
}

// Token returns the token for the given node
func (a *AST) Token(id NodeID) token.Token {
	return a.node[id].token()
}

// Token returns the token for the given node
func (a *AST) Kind(id NodeID) Kind {
	return a.node[id].kind()
}

// Type returns the type of the given node
func (a *AST) Type(id NodeID) types.Type {
	if id >= NodeID(len(a.typ)) {
		return nil
	}
	return a.typ[id]
}

// SetType sets the type of the given node
func (a *AST) SetType(id NodeID, typ types.Type) {
	for id >= NodeID(len(a.typ)) {
		a.typ = append(a.typ, nil)
	}
	a.typ[id] = typ
}

// NumChildren returns the number of children for the given node
func (a *AST) NumChildren(id NodeID) int {
	start := a.node[id].firstChild()
	end := len(a.child)
	if id+1 < NodeID(len(a.node)) {
		end = a.node[id+1].firstChild()
	}
	return int(end - start)
}

// Child returns the nth child of the given node
func (a *AST) Child(id NodeID, index int) NodeID {
	start := a.node[id].firstChild()
	end := len(a.child)
	if id+1 < NodeID(len(a.node)) {
		end = a.node[id+1].firstChild()
	}
	if index < 0 || index >= int(end-start) {
		return InvalidNode
	}
	return a.child[start+index]
}

// Child returns the children of a node
func (a *AST) Children(id NodeID) []NodeID {
	start := a.node[id].firstChild()
	end := len(a.child)
	if id+1 < NodeID(len(a.node)) {
		end = a.node[id+1].firstChild()
	}

	return a.child[start:end]
}

// AddNode adds a new node to the AST
func (a *AST) AddNode(kind Kind, token token.Token, children ...NodeID) NodeID {
	id := NodeID(len(a.node))
	firstchild := len(a.child)
	a.node = append(a.node, newNode(kind, token, firstchild))
	a.child = append(a.child, children...)
	return id
}

// PositionOf returns the line and column of the given token
func (a *AST) PositionOf(tok token.Token) (line int, col int) {
	offset := tok.Offset()
	var lineoffset int
	for i, ch := range a.Src {
		if i >= offset {
			break
		}
		if ch == '\n' {
			line++
			lineoffset = i
		}
	}
	col = offset - lineoffset
	return line + 1, col + 1
}

func (a *AST) SymbolOf(n NodeID) *types.Symbol {
	switch a.node[n].kind() {
	case Name:
		name := a.NodeString(n)
		sym := a.symtab.Lookup(name)
		if sym.Offset == -1 {
			sym.Offset = a.nextAddr
			a.nextAddr++
		}
		return sym
	}
	return nil
}

// NodeBytes returns a cheap byte slice of the node text
func (a *AST) NodeBytes(id NodeID) []byte {
	t := a.node[id].token()
	return a.Src[t.Offset():t.EndOfToken(a.Src)]
}

// NodeString allocates a new string copy of the node text
func (a *AST) NodeString(id NodeID) string {
	return string(a.NodeBytes(id))
}

// String returns a string representation of the AST
func (a *AST) String() string {
	return a.nodeString(a.Root(), "")
}

// StringOf returns a string representation of the AST starting at the given node
func (a *AST) StringOf(root NodeID) string {
	return a.nodeString(root, "")
}

func (a *AST) nodeString(id NodeID, indent string) string {
	if id == 0 {
		return indent + "nil"
	}

	kind := a.Kind(id)

	childnodes := a.Children(id)
	if a.Kind(id).IsTerminal() {
		return fmt.Sprintf("%s%s(%q)", indent, kind.String(), a.NodeBytes(id))
	}

	if len(childnodes) == 0 {
		toktext := ""
		if kind.UsesToken() {
			toktext = fmt.Sprintf("%q", a.NodeBytes(id))
		}

		return fmt.Sprintf("%s%s(%s)", indent, kind.String(), toktext)
	}

	if len(childnodes) <= 2 {
		allterm := true
		for i := range childnodes {
			if !a.Kind(childnodes[i]).IsTerminal() {
				allterm = false
				break
			}
		}
		if allterm {
			str := fmt.Sprintf("%s%s(", indent, kind.String())
			if kind.UsesToken() {
				str += fmt.Sprintf("%q, ", a.NodeBytes(id))
			}

			for i := range childnodes {
				str += a.nodeString(childnodes[i], "")
				if i < len(childnodes)-1 {
					str += ", "
				}
			}
			return str + ")"
		}
	}

	oper := ""
	if kind.UsesToken() {
		oper = fmt.Sprintf("%q, ", a.NodeBytes(id))
	}
	str := fmt.Sprintf("%s%s(%s\n", indent, kind.String(), oper)
	for i := range childnodes {
		str += a.nodeString(childnodes[i], indent+"\t") + ",\n"
	}

	return str + indent + ")"
}

func (a *AST) StackSize() int {
	return a.nextAddr
}
