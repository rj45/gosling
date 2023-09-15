package token

type Kind uint8

const (
	Illegal Kind = iota
	EOF

	Int

	Add
	Sub
	Mul
	Div

	LParen
	RParen

	NumTokens
)

var kindStrs = [...]string{
	Illegal: "Illegal",
	EOF:     "EOF",
	Int:     "Int",
	Add:     "Add",
	Sub:     "Sub",
	Mul:     "Mul",
	Div:     "Div",
	LParen:  "LParen",
	RParen:  "RParen",
}

func (k Kind) String() string {
	return kindStrs[k]
}
