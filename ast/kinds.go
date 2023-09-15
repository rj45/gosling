package ast

// Kind is the kind of AST node
type Kind uint8

const (
	IllegalNode Kind = iota

	Literal
	BinaryExpr
)

var kindNames = []string{
	IllegalNode: "IllegalNode",
	Literal:     "Literal",
	BinaryExpr:  "BinaryExpr",
}

func (k Kind) String() string {
	if k >= Kind(len(kindNames)) {
		return "Unknown"
	}
	return kindNames[k]
}

// IsTerminal returns true if the node is a terminal (leaf) node
func (k Kind) IsTerminal() bool {
	return k == Literal
}

// UsesToken returns true if the node uses a token for its value
func (k Kind) UsesToken() bool {
	return k == BinaryExpr
}
