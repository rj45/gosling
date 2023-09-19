package vm

import (
	"fmt"
	"strconv"
)

// Asm is a simple assembler for the virtual machine.
// It is used by the code generator to generate the
// virtual machine instructions.
type Asm struct {
	Program []Instr

	labels map[string]int
	refs   map[string][]int
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

func (a *Asm) Prologue(locals int) {
	a.instr1(Prologue, locals)
}

func (a *Asm) WordSize() int {
	return 1
}

func (a *Asm) Push() {
	a.instr(Push)
}

func (a *Asm) Pop() {
	a.instr(Pop)
}

func (a *Asm) LoadLocal(offset int) {
	a.instr1(LoadLocal, offset)
}

func (a *Asm) Load() {
	a.instr(Load)
}

func (a *Asm) Store() {
	a.instr(Store)
}

func (a *Asm) LoadInt(lit string) {
	val, err := strconv.Atoi(lit)
	if err != nil {
		panic(err)
	}
	a.instr1(LoadInt, val)
}

func (a *Asm) LocalAddr(offset int) {
	a.instr1(LocalAddr, offset)
}

func (a *Asm) Add() {
	a.instr(Add)
}

func (a *Asm) Sub() {
	a.instr(Sub)
}

func (a *Asm) Mul() {
	a.instr(Mul)
}

func (a *Asm) Div() {
	a.instr(Div)
}

func (a *Asm) Neg() {
	a.instr(Neg)
}

func (a *Asm) Eq() {
	a.instr(Eq)
}

func (a *Asm) Ne() {
	a.instr(Ne)
}

func (a *Asm) Lt() {
	a.instr(Lt)
}

func (a *Asm) Le() {
	a.instr(Le)
}

func (a *Asm) Gt() {
	a.instr(Gt)
}

func (a *Asm) Ge() {
	a.instr(Ge)
}

func (a *Asm) JumpToEpilogue() {
	a.jump(Jump, "epilogue0")
}

func (a *Asm) JumpIfFalse(label string, offset int) {
	name := fmt.Sprintf("%s%d", label, offset)
	a.jump(JumpIfFalse, name)
}

func (a *Asm) Jump(label string, offset int) {
	name := fmt.Sprintf("%s%d", label, offset)
	a.jump(Jump, name)
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

func (a *Asm) Label(label string, offset int) {
	name := fmt.Sprintf("%s%d", label, offset)
	loc := len(a.Program)

	// fixup any references to this label
	if refs, found := a.refs[name]; found {
		for _, ref := range refs {
			a.Program[ref] |= Instr(loc) << 8
		}
		delete(a.refs, name)
	}

	if a.labels == nil {
		a.labels = make(map[string]int)
	}

	// record the label for future references
	a.labels[name] = loc
}

func (a *Asm) Epilogue() {
	a.Label("epilogue", 0)
	a.instr(Epilogue)
}
