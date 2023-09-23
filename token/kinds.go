package token

type Kind uint8

const (
	Illegal Kind = iota
	EOF

	Semicolon
	Comma

	Ident
	Int

	Assign // =
	Define // :=

	Add
	Sub
	Div
	Star
	And

	Eq
	Ne
	Lt
	Le
	Gt
	Ge

	LParen
	RParen
	LBrace
	RBrace

	// keywords
	Return
	If
	Else
	For

	NumTokens
)

var kindStrs = [...]string{
	Illegal:   "Illegal",
	EOF:       "EOF",
	Semicolon: "Semicolon",
	Comma:     "Comma",
	Ident:     "Ident",
	Int:       "Int",
	Assign:    "Assign",
	Add:       "Add",
	Sub:       "Sub",
	Star:      "Mul",
	Div:       "Div",
	Eq:        "Eq",
	Ne:        "Ne",
	Lt:        "Lt",
	Le:        "Le",
	Gt:        "Gt",
	Ge:        "Ge",
	LParen:    "LParen",
	RParen:    "RParen",
	LBrace:    "LBrace",
	RBrace:    "RBrace",
	Return:    "Return",
	If:        "If",
	Else:      "Else",
	For:       "For",
}

func (k Kind) String() string {
	return kindStrs[k]
}
