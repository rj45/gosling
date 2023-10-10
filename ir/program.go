package ir

import (
	"io"
	"strings"

	"github.com/rj45/gosling/types"
)

// Program represents a Go program, and is the root of the IR.
type Program struct {
	*File
	fn    []*Func
	types *types.Universe
}

// NewProgram creates a new Program.
func NewProgram(file *File) *Program {
	return &Program{File: file, types: types.NewUniverse()}
}

// NewFunc creates a new function in the program.
func (p *Program) NewFunc() *Func {
	fn := NewFunc(p.File)
	p.fn = append(p.fn, fn)
	return fn
}

// NumFuncs returns the number of functions in the program.
func (p *Program) NumFuncs() int {
	return len(p.fn)
}

// Func returns the function at the given index.
func (p *Program) Func(index int) *Func {
	return p.fn[index]
}

// FuncNamed returns the function with the given name.
func (p *Program) FuncNamed(name string) *Func {
	for _, fn := range p.fn {
		if fn.Name == name {
			return fn
		}
	}
	return nil
}

// Types returns the universe of types in the program.
func (p *Program) Types() *types.Universe {
	return p.types
}

// SetTypes sets the universe of types in the program.
func (p *Program) SetTypes(uni *types.Universe) {
	p.types = uni
}

// String returns the filename.
func (p *Program) String() string {
	return p.File.String()
}

func (p *Program) Dump() string {
	w := &strings.Builder{}
	p.dump(w)
	return w.String()
}

func (p *Program) dump(w io.Writer) {
	for _, fn := range p.fn {
		fn.dump(w)
	}
}
