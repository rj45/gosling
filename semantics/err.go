package semantics

import "github.com/rj45/gosling/token"

// SemErr is a semantic error, designed to be
// easy to convert to a parser.ParseErr.
type SemErr struct {
	Token token.Token
	Msg   string
}
