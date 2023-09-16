package token

type Kind uint8

const (
	Illegal Kind = iota
	EOF
	Semicolon

	Ident
	Int

	Assign

	Add
	Sub
	Mul
	Div

	Eq
	Ne
	Lt
	Le
	Gt
	Ge

	LParen
	RParen

	NumTokens
)

var kindStrs = [...]string{
	Illegal:   "Illegal",
	EOF:       "EOF",
	Semicolon: "Semicolon",
	Ident:     "Ident",
	Int:       "Int",
	Assign:    "Assign",
	Add:       "Add",
	Sub:       "Sub",
	Mul:       "Mul",
	Div:       "Div",
	Eq:        "Eq",
	Ne:        "Ne",
	Lt:        "Lt",
	Le:        "Le",
	Gt:        "Gt",
	Ge:        "Ge",
	LParen:    "LParen",
	RParen:    "RParen",
}

func (k Kind) String() string {
	return kindStrs[k]
}