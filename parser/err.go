package parser

import (
	"bytes"
	"fmt"

	"github.com/rj45/gosling/token"
)

type ParseErr struct {
	p   *Parser
	tok token.Token
	msg string
}

func (e ParseErr) Error() string {
	buf := &bytes.Buffer{}
	line, col := e.p.ast.PositionOf(e.tok)
	fmt.Fprintf(buf, "error %d:%d: %s\n", line, col, e.msg)

	// print a few lines before and after the error
	lines := bytes.Split(e.p.src, []byte("\n"))
	for i := line - 3; i <= line+3; i++ {
		if i < 0 || i >= len(lines) {
			continue
		}
		fmt.Fprintf(buf, "%s\n", lines[i])
		if i == line-1 {
			fmt.Fprintf(buf, "%*s^ here\n", col-1, "")
		}
	}
	return buf.String()
}

var tokenErrStrs = map[token.Kind]string{
	token.EOF:       "end of file",
	token.Semicolon: "newline or ';' or '}'",
	token.RBrace:    "'}'",
}
