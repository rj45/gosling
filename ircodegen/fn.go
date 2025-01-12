package ircodegen

import (
	"github.com/rj45/gosling/ir"
)

func (g *CodeGen) genFunction() {
	// TODO: determine stack size
	g.asm.Prologue(g.fn.Name, 0)
	g.visited = make([]bool, g.fn.Graph.NumNodes()+1)

	// g.fn.Graph.Schedule()

	// g.fn.Graph.IterBlocks()(func(block int, node ir.Node) bool {
	// 	if node.Op().Opcode() == ir.OpReturn {
	// 		g.genReturn(node)
	// 	}
	// 	return true
	// })
	var last ir.Node
	g.fn.Graph.DCEIter()(func(node ir.Node) bool {
		last = node
		return true
	})

	g.gen(last)

	g.asm.Epilogue()
}

func (g *CodeGen) genBinaryExpr(node ir.Node) {
	g.gen(node.Input(0))
	g.asm.Push()
	g.gen(node.Input(1))
	g.asm.Pop(1)

	switch node.Op().Opcode() {
	case ir.OpAdd:
		g.asm.Add()
	case ir.OpSub:
		g.asm.Sub()
	case ir.OpMul:
		g.asm.Mul()
	case ir.OpDiv:
		g.asm.Div()
	default:
		panic("unknown binary expr op " + node.Op().String())
	}
}

func (g *CodeGen) gen(node ir.Node) {
	if g.visited[node.Index()] {
		return
	}
	g.visited[node.Index()] = true

	switch node.Op().Opcode() {
	case ir.OpStart:
		// do nothing
	case ir.OpConst:
		value := g.fn.Types.StringOf(node.Type())
		if value == "false" {
			g.asm.LoadInt("0")
		} else if value == "true" {
			g.asm.LoadInt("1")
		} else {
			g.asm.LoadInt(g.fn.Types.StringOf(node.Type()))
		}

	case ir.OpReturn:
		node.InputsIter()(func(_ int, input ir.Node) bool {
			g.gen(input)
			return true
		})
		g.asm.JumpToEpilogue()
	case ir.OpAdd, ir.OpSub, ir.OpMul, ir.OpDiv:
		g.genBinaryExpr(node)
	case ir.OpNeg:
		g.gen(node.Input(0))
		g.asm.Neg()

	case ir.OpIf:
		g.gen(node.Input(0))
		g.gen(node.Input(1))

	case ir.OpThen:
		g.gen(node.Input(0))
	case ir.OpElse:
		g.gen(node.Input(0))
	case ir.OpRegion:
		// find all the phi nodes, and extract their values
		var thens []ir.Node
		var elses []ir.Node
		node.UseIter()(func(use ir.Use) bool {
			user := use.User()
			if user.Op().Opcode() == ir.OpPhi {
				thens = append(thens, user.Input(1))
				elses = append(elses, user.Input(2))
			}
			return true
		})

		// generate the then clause
		g.gen(node.Input(0))
		for _, then := range thens {
			g.gen(then)
		}

		// jump over any else clauses
		// TODO: detect an empty else and avoid the jump
		g.asm.Jump(node.Name(), g.label)

		// generate the else clause
		g.gen(node.Input(1))
		for _, els := range elses {
			g.gen(els)
		}

		// generate the merge point
		g.asm.Label(node.Name(), g.label)

	case ir.OpStop:
		// generate returns starting at the last input
		for i := node.NumInputs() - 1; i >= 0; i-- {
			g.gen(node.Input(i))
		}

	default:
		panic("unknown node op " + node.Op().String())
	}
}
