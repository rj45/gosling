package hlir

import (
	"strconv"

	"github.com/rj45/gosling/ir"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

type ref struct {
	instr ir.ValueID // the instruction to fixup
	oper  int        // the operand to fixup
}

type Assembler struct {
	*ir.Program

	// The current function being assembled.
	// This is set by the Prologue and Epilogue operators.
	Func *ir.Func

	// The current basic block being assembled.
	// This is set by the Label operator.
	Block *ir.Block

	a ir.ValueID
	b ir.ValueID

	labels map[string]ir.ValueID
	refs   map[string][]ref
}

func NewAssembler(file *token.File) *Assembler {
	return &Assembler{
		Program: ir.NewProgram(file),
	}
}

func (a *Assembler) WordSize() int {
	return 1
}

func (a *Assembler) Types(uni *types.Universe) {
	a.Program.SetTypes(uni)
}

func (a *Assembler) DeclareFunction(fnname string, sig types.Type) {
	fn := a.Program.NewFunc()
	fn.Name = fnname
	fn.Sig = sig
}

func (a *Assembler) Prologue(name string, numLocals int) {
	a.Func = a.Program.FuncNamed(name)

	a.Label(name+".entry", 0)

	a.Block.AddValue(Prologue, 0, types.Void, a.Func.ValueForConst(ir.IntConst(int64(numLocals))))

	a.a = ir.InvalidValue
	a.b = ir.InvalidValue
}

func (a *Assembler) Epilogue() {
	a.Label(a.Func.Name+".epilogue", 0)
	a.Block.AddValue(Epilogue, 0, types.Void)
	ft := a.Program.Types().Func(a.Func.Sig)
	if ft.ReturnType() == types.Void {
		a.Block.UpdateTerminator(Return)
		return
	}

	a.Block.UpdateTerminator(Return, a.a)

	a.Block = nil
	a.Func = nil
}

func (a *Assembler) Push() {
	a.Block.AddValue(Push, 0, types.Void, a.a)
}

func (a *Assembler) Pop(reg int) {
	a.b = a.Block.AddValue(Pop, 0, types.Int)
	a.Func.AddReg(a.b, ir.RegID(reg))
}

func (a *Assembler) LoadLocal(index int) {
	a.a = a.Block.AddValue(LoadLocal, 0, types.Int, a.Func.ValueForConst(ir.IntConst(int64(index))))
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) StoreLocal(index int) {
	reg := a.Func.ValueForConst(ir.RegConst(ir.NewRegMask(ir.RegID(index))))
	c := a.Func.ValueForConst(ir.IntConst(int64(index)))
	a.Block.AddValue(StoreLocal, 0, types.Void, reg, c)
}

func (a *Assembler) LoadInt(value string) {
	ival, _ := strconv.ParseInt(value, 10, 64)
	a.a = a.Block.AddValue(LoadInt, 0, types.Int, a.Func.ValueForConst(ir.IntConst(ival)))
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Load() {
	a.a = a.Block.AddValue(Load, 0, types.Int, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Store() {
	a.Block.AddValue(Store, 0, types.Void, a.a, a.b)
}

func (a *Assembler) LocalAddr(index int) {
	a.a = a.Block.AddValue(LocalAddr, 0, types.Int, a.Func.ValueForConst(ir.IntConst(int64(index))))
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Add() {
	a.a = a.Block.AddValue(Add, 0, types.Int, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Sub() {
	a.a = a.Block.AddValue(Sub, 0, types.Int, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Mul() {
	a.a = a.Block.AddValue(Mul, 0, types.Int, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Div() {
	a.a = a.Block.AddValue(Div, 0, types.Int, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Neg() {
	a.a = a.Block.AddValue(Neg, 0, types.Int, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Eq() {
	a.a = a.Block.AddValue(Eq, 0, types.Bool, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Ne() {
	a.a = a.Block.AddValue(Ne, 0, types.Bool, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Lt() {
	a.a = a.Block.AddValue(Lt, 0, types.Bool, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Gt() {
	a.a = a.Block.AddValue(Gt, 0, types.Bool, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Le() {
	a.a = a.Block.AddValue(Le, 0, types.Bool, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Ge() {
	a.a = a.Block.AddValue(Ge, 0, types.Bool, a.b, a.a)
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Call(fnname string) {
	fn := a.Program.FuncNamed(fnname)
	rettype := a.Program.Types().Func(fn.Sig).ReturnType()
	a.a = a.Block.AddValue(Call, 0, rettype, a.Func.ValueForConst(ir.FuncConst(fn)))
	a.Func.AddReg(a.a, ir.R0)
}

func (a *Assembler) Jump(label string, id int) {
	a.jump(Jump, label, id)
}

func (a *Assembler) jump(op Op, label string, id int, ops ...ir.ValueID) {
	name := label + strconv.Itoa(id)
	if bid, found := a.labels[name]; found {
		a.Block.UpdateTerminator(op, append(ops, bid)...)
		return
	}

	a.Block.UpdateTerminator(op, append(ops, ir.InvalidValue)...)

	if a.refs == nil {
		a.refs = make(map[string][]ref)
	}
	a.refs[name] = append(a.refs[name], ref{a.Block.Terminator(), len(ops)})
}

func (a *Assembler) JumpIfFalse(label string, id int) {
	a.jump(If, label, id, a.a, ir.InvalidValue)
}

func (a *Assembler) JumpToEpilogue() {
	a.jump(Jump, a.Func.Name+".epilogue", 0)
}

func (a *Assembler) Label(label string, id int) {
	name := label + strconv.Itoa(id)
	bid := a.Func.NewBlock(name, Jump, 0, ir.InvalidValue)

	if a.Block != nil {
		term := a.Block.Terminator()
		if a.Func.Op(term) == If {
			ops := a.Func.Operands(term)
			ops[1] = bid
			a.Func.SetOperands(term, ops...)
		} else {
			lastOper := a.Func.NumOperands(term) - 1
			if a.Func.Operand(term, lastOper) == ir.InvalidValue {
				ops := a.Func.Operands(a.Block.Terminator())
				ops[lastOper] = bid
				a.Func.SetOperands(a.Block.Terminator(), ops...)
			}
		}
	}

	a.Block = a.Func.Block(bid)

	// fixup any references to this label
	if refs, found := a.refs[name]; found {
		for _, ref := range refs {
			a.Func.SetOperand(ref.instr, ref.oper, bid)
		}
		delete(a.refs, name)
	}

	if a.labels == nil {
		a.labels = make(map[string]ir.ValueID)
	}

	// record the label for future references
	a.labels[name] = bid
}
