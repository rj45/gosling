package vm

import (
	"log"

	"github.com/rj45/gosling/ir"
)

type Asm struct {
	Program []Instr

	labels map[string]int
	refs   map[string][]int
	fn     string
}

func NewAsm() *Asm {
	return &Asm{}
}

func (a *Asm) instr(op Opcode) {
	a.Program = append(a.Program, Instr(op))
}

func (a *Asm) instr1(op Opcode, arg int) {
	a.Program = append(a.Program, Instr(op)|Instr(arg)<<8)
}

func (a *Asm) Prologue(fn string, locals int) {
	a.fn = "_" + fn
	a.Label(a.fn)
	a.instr1(Prologue, locals)
}

func (a *Asm) Epilogue() {
}

func (a *Asm) Push(src ir.RegMask) {
	if !src.HasReg(ir.R0) {
		panic("src must be R0")
	}
	a.instr(Push)
}

func (a *Asm) Pop(dest ir.RegMask) {
	if len(dest.Regs()) != 1 {
		panic("dest must have one register")
	}
	a.instr1(Pop, int(dest.Pop()))
}

func (a *Asm) LoadLocal(dest ir.RegMask, local int) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	a.instr1(LoadLocal, local)
}

func (a *Asm) StoreLocal(dest ir.RegMask, local int) {
	if !dest.HasReg(ir.RegID(local)) {
		log.Panicf("dest must be r%d, was: %s", local, dest)
	}
	a.instr1(StoreLocal, local)
}

func (a *Asm) Load(dest ir.RegMask, src ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src.HasReg(ir.R0) {
		panic("src must be R0")
	}
	a.instr(Load)
}

func (a *Asm) Store(src ir.RegMask, addr ir.RegMask) {
	if !src.HasReg(ir.R0) {
		panic("src must be R0")
	}
	if !addr.HasReg(ir.R1) {
		panic("addr must be R1")
	}
	a.instr(Store)
}

func (a *Asm) LoadInt(dest ir.RegMask, imm int64) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	a.instr1(LoadInt, int(imm))
}

func (a *Asm) LocalAddr(dest ir.RegMask, local int) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	a.instr1(LocalAddr, local)
}

func (a *Asm) Add(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Add)
}

func (a *Asm) Sub(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Sub)
}

func (a *Asm) Mul(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Mul)
}

func (a *Asm) Div(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Div)
}

func (a *Asm) Neg(dest ir.RegMask, src ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src.HasReg(ir.R0) {
		panic("src must be R0")
	}
	a.instr(Neg)
}

func (a *Asm) Eq(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Eq)
}

func (a *Asm) Ne(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Ne)
}

func (a *Asm) Lt(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Lt)
}

func (a *Asm) Le(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Le)
}

func (a *Asm) Gt(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Gt)
}

func (a *Asm) Ge(dest ir.RegMask, src1 ir.RegMask, src2 ir.RegMask) {
	if !dest.HasReg(ir.R0) {
		panic("dest must be R0")
	}
	if !src1.HasReg(ir.R1) {
		panic("src1 must be R1")
	}
	if !src2.HasReg(ir.R0) {
		panic("src2 must be R0")
	}
	a.instr(Ge)
}

func (a *Asm) Call(fn string) {
	a.jump(Call, "_"+fn)
}

func (a *Asm) If(test ir.RegMask, then string, els string) {
	if !test.HasReg(ir.R0) {
		panic("test must be R0")
	}
	a.jump(JumpIfFalse, els)
	a.jump(Jump, then)
}

func (a *Asm) Jump(label string) {
	a.jump(Jump, label)
}

func (a *Asm) jump(op Opcode, label string) {
	if loc, found := a.labels[label]; found {
		a.instr1(op, loc)
		return
	}
	if a.refs == nil {
		a.refs = make(map[string][]int)
	}
	a.refs[label] = append(a.refs[label], len(a.Program))
	a.instr1(op, 0)
}

func (a *Asm) Label(label string) {
	loc := len(a.Program)

	// fixup any references to this label
	if refs, found := a.refs[label]; found {
		for _, ref := range refs {
			a.Program[ref] |= Instr(loc) << 8
		}
		delete(a.refs, label)
	}

	if a.labels == nil {
		a.labels = make(map[string]int)
	}

	// record the label for future references
	a.labels[label] = loc
}

func (a *Asm) Return() {
	if a.fn == "_main" {
		a.instr(Exit)
		return
	}
	a.instr(Return)
}
