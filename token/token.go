package token

import (
	"strconv"
)

// 7 bits kind (127)
// 25 bits offset (32 MB)
type Token uint32

type Kind uint8

const (
	Illegal Kind = iota
	EOF

	Int

	Add
	Sub

	NumTokens
)

// Kind returns the TokenKind
func (t Token) Kind() Kind {
	// least significant 7 bits
	return Kind(t & 0x7f)
}

// Offset returns the offset into the src []byte
func (t Token) Offset() int {
	// remaining 19 bits
	return int(t >> 7)
}

func NewToken(kind Kind, offset int) Token {
	return Token(uint32(offset)<<7 | uint32(kind))
}

func (t Token) EndOfToken(src []byte) int {
	// eot: End Of Token, starts with the current token's offset
	eot := t.Offset()

	// Jump to the end of the current token
	switch t.Kind() {
	case Int:
		for eot < len(src) && src[eot] >= '0' && src[eot] <= '9' {
			eot++
		}

	case Illegal, EOF:
		// zero length

	case Add, Sub:
		eot++ // For single character tokens (like '+', '-', etc.)
	default:
		panic("todo: handle other tokens")
	}

	return eot
}

// Next returns the next Token in src relative to the current Token.
func (t Token) Next(src []byte) Token {
	pos := t.EndOfToken(src) // Set the start position for scanning the next token

	// Skip over whitespace and comments
skip:
	for pos < len(src) {
		switch src[pos] {
		case '\n':
			pos++
		case ' ', '\t', '\r':
			pos++
		case '/':
			if pos+1 < len(src) && src[pos+1] == '/' {
				// Skip entire line for line comments
				pos += 2
				for pos < len(src) && src[pos] != '\n' {
					pos++
				}
			} else if pos+1 < len(src) && src[pos+1] == '*' {
				// Skip until closing for block comments
				pos += 2
				for pos+1 < len(src) && !(src[pos] == '*' && src[pos+1] == '/') {
					pos++
				}
				pos += 2
			} else {
				break skip
			}
		default:
			break skip
		}
	}

	if pos >= len(src) {
		return NewToken(EOF, len(src))
	}

	switch ch := src[pos]; {
	case ch == '+':
		// Handle other cases like '+=', '++' ...
		return NewToken(Add, pos)
	case ch == '-':
		// Handle other cases like '-=', '--' ...
		return NewToken(Sub, pos)
	case ch >= '0' && ch <= '9':
		start := pos
		for pos < len(src) && src[pos] >= '0' && src[pos] <= '9' {
			pos++
		}
		return NewToken(Int, start)
	default:
		return NewToken(Illegal, pos)
	}
}

// Bytes returns a cheap byte slice of the token text
func (t Token) Bytes(src []byte) []byte {
	return src[t.Offset():t.EndOfToken(src)]
}

// StringFrom allocates a new string copy of the token text
func (t Token) StringFrom(src []byte) string {
	return string(t.Bytes(src))
}

func (t Token) String() string {
	return t.Kind().String() + "(" + strconv.Itoa(t.Offset()) + ")"
}

var kindStrs = [...]string{
	Illegal: "Illegal",
	EOF:     "EOF",
	Int:     "Int",
	Add:     "Add",
	Sub:     "Sub",
}

func (k Kind) String() string {
	return kindStrs[k]
}