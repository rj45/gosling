package ast

// Kind is the kind of AST node
type Kind uint8

const (
	IllegalNode Kind = iota

	Literal
	Name

	BinaryExpr
	UnaryExpr

	StmtList
	ExprStmt
	AssignStmt
)

var kindNames = []string{
	IllegalNode: "IllegalNode",
	Literal:     "Literal",
	Name:        "Name",
	BinaryExpr:  "BinaryExpr",
	UnaryExpr:   "UnaryExpr",
	StmtList:    "StmtList",
	ExprStmt:    "ExprStmt",
	AssignStmt:  "AssignStmt",
}

func (k Kind) String() string {
	if k >= Kind(len(kindNames)) {
		return "Unknown"
	}
	return kindNames[k]
}

// IsTerminal returns true if the node is a terminal (leaf) node
func (k Kind) IsTerminal() bool {
	return k == Literal || k == Name
}

// UsesToken returns true if the node uses a token for its value
func (k Kind) UsesToken() bool {
	return k == BinaryExpr || k == UnaryExpr
}
