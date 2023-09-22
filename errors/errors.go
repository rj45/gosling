package errors

import (
	"bytes"
	"fmt"

	"github.com/rj45/gosling/token"
)

type Err struct {
	src []byte
	tok token.Token
	msg string
}

func New(src []byte, tok token.Token, msg string) *Err {
	return &Err{src: src, tok: tok, msg: msg}
}

func Newf(src []byte, tok token.Token, msg string, args ...interface{}) *Err {
	return &Err{src: src, tok: tok, msg: fmt.Sprintf(msg, args...)}
}

func (e *Err) Error() string {
	buf := &bytes.Buffer{}
	line, col := e.positionOf(e.tok)
	fmt.Fprintf(buf, "error %d:%d: %s\n", line, col, e.msg)

	// print a few lines before and after the error
	lines := bytes.Split(e.src, []byte("\n"))
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

// PositionOf returns the line and column of the given token
func (e *Err) positionOf(tok token.Token) (line int, col int) {
	offset := tok.Offset()
	var lineoffset int
	for i, ch := range e.src {
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
