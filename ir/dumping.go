package ir

import (
	"fmt"
	"io"
	"strings"

	"github.com/rj45/gosling/types"
)

func (p *Program) Dump() string {
	w := &strings.Builder{}
	p.dump(w)
	return w.String()
}

func (p *Program) dump(w io.Writer) {
	for _, pkg := range p.Packages {
		pkg.dump(w)
	}
}

func (p *Package) dump(w io.Writer) {
	for _, fn := range p.Funcs {
		fn.dump(w)
	}
}

// Dump returns a string representation of the function.
func (fn *Function) Dump() string {
	w := &strings.Builder{}
	fn.dump(w)
	return w.String()
}

func (fn *Function) dump(w io.Writer) {
	sig := fn.Types.Func(fn.Type)
	params := make([]string, len(sig.ParamTypes()))
	for i, typ := range sig.ParamTypes() {
		params[i] = fn.Types.StringOf(typ)
	}
	pstr := "(" + strings.Join(params, ", ") + ")"
	if sig.ReturnType() != types.Void {
		pstr += " " + fn.Types.StringOf(sig.ReturnType())
	}
	fmt.Fprintln(w, "func", fn.Name+pstr, "{")
	fn.Graph.dump(w, "\t", fn.Types)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}

func (g *Graph) dump(w io.Writer, indent string, types *types.Universe) {
	// TODO: do a post-order traversal, and then dump in reverse order
	for i := range g.nodes {
		g.dumpNode(w, NodeID(i), indent, types)
	}
}

func (g *Graph) dumpNode(w io.Writer, id NodeID, indent string, types *types.Universe) {
	name := g.names[id]
	namestr := name.name
	if g.lastNameNum[name.name] > 1 {
		namestr += fmt.Sprintf("%d", name.num)
	}

	if g.op(id).IsControlFlow() && g.op(id).Opcode() != OpStart && g.op(id).Opcode() != OpReturn {
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "%s$%s = %s", indent, namestr, g.op(id).String())
	for i := 0; i < g.numInputs(id); i++ {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		name := g.names[g.input(id, i).ID()]
		namestr := name.name
		if g.lastNameNum[name.name] > 1 {
			namestr += fmt.Sprintf("%d", name.num)
		}
		fmt.Fprintf(w, " $%s", namestr)
	}
	if g.op(id).Opcode() == OpConst {
		fmt.Fprintf(w, " %s", types.StringOf(g.types[id]))
	}
	fmt.Fprintln(w)
}
