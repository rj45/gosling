package token

import (
	"strconv"
)

// 7 bits kind (127)
// 25 bits offset (32 MB)
type Token uint32

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
	case Ident:
		for eot < len(src) && (src[eot] >= 'a' && src[eot] <= 'z' || src[eot] >= 'A' && src[eot] <= 'Z' || src[eot] >= '0' && src[eot] <= '9' || src[eot] == '_') {
			eot++
		}

	case Return, If, Else, For:
		// for keywords, assume kind length is the token length
		eot += len(t.Kind().String())

	case Illegal, EOF:
		// zero length

	case Add, Sub, And, Star, Div, LParen, RParen, LBrace, RBrace, Lt, Gt, Semicolon, Assign:
		eot++ // For single character tokens (like '+', '-', etc.)
	case Eq, Ne, Le, Ge:
		eot += 2 // For double character tokens (like '==', '!=', etc.)
	default:
		panic("todo: handle other tokens")
	}

	return eot
}

// nlsemi is a list of tokens that have the a following newline converted to a semicolon
var nlsemi = [NumTokens]bool{
	Int:    true,
	Ident:  true,
	RParen: true,
	RBrace: true,
	Return: true,
}

var keywords = map[string]Kind{
	"return": Return,
	"if":     If,
	"else":   Else,
	"for":    For,
}

// Next returns the next Token in src relative to the current Token.
func (t Token) Next(src []byte) Token {
	pos := t.EndOfToken(src)   // Set the start position for scanning the next token
	nlsemi := nlsemi[t.Kind()] // Whether or not to treat a newline as a semicolon

	// Skip over whitespace and comments
skip:
	for pos < len(src) {
		switch src[pos] {
		case '\n':
			if nlsemi {
				return NewToken(Semicolon, pos)
			}
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
	case ch == ';':
		return NewToken(Semicolon, pos)
	case ch == '+':
		return NewToken(Add, pos)
	case ch == '-':
		return NewToken(Sub, pos)
	case ch == '*':
		return NewToken(Star, pos)
	case ch == '/':
		return NewToken(Div, pos)
	case ch == '&':
		return NewToken(And, pos)

	case ch == '(':
		return NewToken(LParen, pos)
	case ch == ')':
		return NewToken(RParen, pos)
	case ch == '{':
		return NewToken(LBrace, pos)
	case ch == '}':
		return NewToken(RBrace, pos)

	case ch == '=':
		if pos+1 < len(src) && src[pos+1] == '=' {
			return NewToken(Eq, pos)
		}
		return NewToken(Assign, pos)
	case ch == '!':
		if pos+1 < len(src) && src[pos+1] == '=' {
			return NewToken(Ne, pos)
		}
	case ch == '<':
		if pos+1 < len(src) && src[pos+1] == '=' {
			return NewToken(Le, pos)
		}
		return NewToken(Lt, pos)
	case ch == '>':
		if pos+1 < len(src) && src[pos+1] == '=' {
			return NewToken(Ge, pos)
		}
		return NewToken(Gt, pos)

	case ch >= '0' && ch <= '9':
		start := pos
		for pos < len(src) && src[pos] >= '0' && src[pos] <= '9' {
			pos++
		}
		if src[pos] >= 'a' && src[pos] <= 'z' || src[pos] >= 'A' && src[pos] <= 'Z' {
			// identifier starting with a number
			return NewToken(Illegal, start)
		}
		return NewToken(Int, start)
	case ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_':
		start := pos
		for pos < len(src) && (src[pos] >= 'a' && src[pos] <= 'z' || src[pos] >= 'A' && src[pos] <= 'Z' || src[pos] >= '0' && src[pos] <= '9' || src[pos] == '_') {
			pos++
		}
		if keyword, ok := keywords[string(src[start:pos])]; ok {
			return NewToken(keyword, start)
		}
		return NewToken(Ident, start)
	}
	return NewToken(Illegal, pos)
}

// NextValidToken returns the next non-illegal token.
func (t Token) NextValidToken(src []byte) Token {
	for t.Kind() == Illegal {
		pos := t.Offset()

		skipped := false

		// skip over alphanumeric characters
		for pos < len(src) && (src[pos] >= 'a' && src[pos] <= 'z' || src[pos] >= 'A' && src[pos] <= 'Z' || src[pos] >= '0' && src[pos] <= '9' || src[pos] == '_') {
			skipped = true
			pos++
		}

		if !skipped && src[pos] != ' ' && src[pos] != '\t' && src[pos] != '\r' && src[pos] != '\n' {
			pos++
		}

		t = NewToken(Illegal, pos)
		t = t.Next(src)
	}
	return t
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
