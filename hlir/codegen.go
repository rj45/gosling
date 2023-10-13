package hlir

import (
	"github.com/rj45/gosling/ir"
)

type Assembler interface {
	Prologue(string, int)
	Epilogue()

	Push(ir.RegMask)
	Pop(ir.RegMask)
	LoadLocal(ir.RegMask, int)
	StoreLocal(ir.RegMask, int)
	Load(ir.RegMask, ir.RegMask)
	Store(ir.RegMask, ir.RegMask)

	LoadInt(ir.RegMask, int64)
	LocalAddr(ir.RegMask, int)

	Add(ir.RegMask, ir.RegMask, ir.RegMask)
	Sub(ir.RegMask, ir.RegMask, ir.RegMask)
	Mul(ir.RegMask, ir.RegMask, ir.RegMask)
	Div(ir.RegMask, ir.RegMask, ir.RegMask)

	Neg(ir.RegMask, ir.RegMask)

	Eq(ir.RegMask, ir.RegMask, ir.RegMask)
	Ne(ir.RegMask, ir.RegMask, ir.RegMask)
	Lt(ir.RegMask, ir.RegMask, ir.RegMask)
	Le(ir.RegMask, ir.RegMask, ir.RegMask)
	Gt(ir.RegMask, ir.RegMask, ir.RegMask)
	Ge(ir.RegMask, ir.RegMask, ir.RegMask)

	Call(string)
	If(ir.RegMask, string, string)
	Jump(string)
	Label(string)
	Return()
}

type CodeGen struct {
	*ir.Program
	asm Assembler

	fn *ir.Func
}

func New(program *ir.Program, asm Assembler) *CodeGen {
	return &CodeGen{
		Program: program,
		asm:     asm,
	}
}

func (c *CodeGen) Generate() {
	for i := 0; i < c.NumFuncs(); i++ {
		c.fn = c.Func(i)
		if c.fn.Name == "main" {
			c.generateFunction()
		}
	}

	for i := 0; i < c.NumFuncs(); i++ {
		c.fn = c.Func(i)
		if c.fn.Name != "main" {
			c.generateFunction()
		}
	}
}

func (c *CodeGen) generateFunction() {
	for i := 0; i < c.fn.NumBlocks(); i++ {
		c.generateBlock(c.fn.BlockAt(i))
	}
}

func (c *CodeGen) generateBlock(blk *ir.Block) {
	c.asm.Label(blk.Name)
	for i := 0; i < blk.NumValues(); i++ {
		c.generateInstr(blk.ValueAt(i))
	}
	c.generateInstr(blk.Terminator())
}

func (c *CodeGen) generateInstr(instr ir.ValueID) {
	reg := [3]ir.RegMask{}
	ri := 0
	if c.fn.Regs(instr) != 0 {
		reg[0] = c.fn.Regs(instr)
		ri = 1
	}
	for i := 0; i < c.fn.NumOperands(instr); i++ {
		if c.fn.Regs(c.fn.Operand(instr, i)) != 0 {
			reg[ri] = c.fn.Regs(c.fn.Operand(instr, i))
			ri++
		}
	}

	switch c.fn.Op(instr) {
	case Prologue:
		v := c.fn.ConstForValue(c.fn.Operand(instr, 0))
		c.asm.Prologue(c.fn.Name, int(v.Value().(int64)))
	case Epilogue:
		c.asm.Epilogue()
	case Push:
		c.asm.Push(reg[0])
	case Pop:
		c.asm.Pop(reg[0])
	case LoadLocal:
		v := c.fn.ConstForValue(c.fn.Operand(instr, 0))
		c.asm.LoadLocal(reg[0], int(v.Value().(int64)))
	case StoreLocal:
		rc := c.fn.ConstForValue(c.fn.Operand(instr, 0))
		reg, ok := ir.RegValue(rc)
		if !ok {
			panic("not a reg: " + rc.String())
		}

		v := c.fn.ConstForValue(c.fn.Operand(instr, 1))
		c.asm.StoreLocal(reg, int(v.Value().(int64)))
	case Load:
		c.asm.Load(reg[0], reg[1])
	case Store:
		c.asm.Store(reg[0], reg[1])
	case LocalAddr:
		v := c.fn.ConstForValue(c.fn.Operand(instr, 0))
		c.asm.LocalAddr(reg[0], int(v.Value().(int64)))
	case LoadInt:
		v := c.fn.ConstForValue(c.fn.Operand(instr, 0))
		c.asm.LoadInt(reg[0], v.Value().(int64))
	case Add:
		c.asm.Add(reg[0], reg[1], reg[2])
	case Sub:
		c.asm.Sub(reg[0], reg[1], reg[2])
	case Mul:
		c.asm.Mul(reg[0], reg[1], reg[2])
	case Div:
		c.asm.Div(reg[0], reg[1], reg[2])
	case Neg:
		c.asm.Neg(reg[0], reg[1])
	case Eq:
		c.asm.Eq(reg[0], reg[1], reg[2])
	case Ne:
		c.asm.Ne(reg[0], reg[1], reg[2])
	case Lt:
		c.asm.Lt(reg[0], reg[1], reg[2])
	case Le:
		c.asm.Le(reg[0], reg[1], reg[2])
	case Gt:
		c.asm.Gt(reg[0], reg[1], reg[2])
	case Ge:
		c.asm.Ge(reg[0], reg[1], reg[2])
	case Call:
		cfn := c.fn.ConstForValue(c.fn.Operand(instr, 0))
		c.asm.Call(cfn.String())
	case Jump:
		b := c.fn.Block(c.fn.Operand(instr, 0))
		dest := b.Name
		c.asm.Jump(dest)
	case If:
		then := c.fn.Block(c.fn.Operand(instr, 1)).Name
		els := c.fn.Block(c.fn.Operand(instr, 2)).Name
		c.asm.If(reg[0], then, els)
	case Return:
		c.asm.Return()
	default:
		panic("unknown op: " + c.fn.Op(instr).String())
	}
}
