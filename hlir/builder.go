package hlir

import (
	"strconv"

	"github.com/rj45/gosling/ir"
	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

type ref struct {
	instr ir.Value // the instruction to fixup
	oper  int      // the operand to fixup
}

type Builder struct {
	*ir.Program

	// The current function being assembled.
	// This is set by the Prologue and Epilogue operators.
	Func *ir.Func

	// The current basic block being assembled.
	// This is set by the Label operator.
	Block *ir.Block

	a ir.Value
	b ir.Value

	labels map[string]ir.Value
	refs   map[string][]ref
}

func NewBuilder(file *token.File) *Builder {
	return &Builder{
		Program: ir.NewProgram(file),
	}
}

func (b *Builder) WordSize() int {
	return 1
}

func (b *Builder) Types(uni *types.Universe) {
	b.Program.SetTypes(uni)
}

func (b *Builder) DeclareFunction(fnname string, sig types.Type) {
	fn := b.Program.NewFunc()
	fn.Name = fnname
	fn.Sig = sig
}

func (b *Builder) Prologue(name string, numLocals int) {
	b.Func = b.Program.FuncNamed(name)

	b.Label(name+".entry", 0)

	b.Block.AddValue(Prologue, 0, types.Void, b.Func.ValueForConst(ir.IntConst(int64(numLocals))))

	b.a = ir.Value{}
	b.b = ir.Value{}
}

func (b *Builder) Epilogue() {
	b.Label(b.Func.Name+".epilogue", 0)
	b.Block.AddValue(Epilogue, 0, types.Void)
	ft := b.Program.Types().Func(b.Func.Sig)
	if ft.ReturnType() == types.Void {
		b.Block.UpdateTerminator(Return)
		return
	}

	b.Block.UpdateTerminator(Return, b.a)

	b.Block = nil
	b.Func = nil
}

func (b *Builder) Push() {
	b.Block.AddValue(Push, 0, types.Void, b.a)
}

func (b *Builder) Pop(reg int) {
	b.b = b.Block.AddValue(Pop, 0, types.Int).AddReg(ir.RegID(reg))
}

func (b *Builder) LoadLocal(index int) {
	b.a = b.Block.AddValueAny(LoadLocal, 0, types.Int, index).AddReg(ir.R0)
}

func (b *Builder) StoreLocal(index int) {
	b.Block.AddValueAny(StoreLocal, 0, types.Void, ir.RegID(index), index)
}

func (b *Builder) LoadInt(value string) {
	ival, _ := strconv.ParseInt(value, 10, 64)
	b.a = b.Block.AddValueAny(LoadInt, 0, types.Int, ival).AddReg(ir.R0)
}

func (b *Builder) Load() {
	b.a = b.Block.AddValue(Load, 0, types.Int, b.a).AddReg(ir.R0)
}

func (b *Builder) Store() {
	b.Block.AddValue(Store, 0, types.Void, b.a, b.b)
}

func (b *Builder) LocalAddr(index int) {
	b.a = b.Block.AddValueAny(LocalAddr, 0, types.Int, index).AddReg(ir.R0)
}

func (b *Builder) Add() {
	b.a = b.Block.AddValue(Add, 0, types.Int, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Sub() {
	b.a = b.Block.AddValue(Sub, 0, types.Int, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Mul() {
	b.a = b.Block.AddValue(Mul, 0, types.Int, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Div() {
	b.a = b.Block.AddValue(Div, 0, types.Int, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Neg() {
	b.a = b.Block.AddValue(Neg, 0, types.Int, b.a).AddReg(ir.R0)
}

func (b *Builder) Eq() {
	b.a = b.Block.AddValue(Eq, 0, types.Bool, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Ne() {
	b.a = b.Block.AddValue(Ne, 0, types.Bool, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Lt() {
	b.a = b.Block.AddValue(Lt, 0, types.Bool, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Gt() {
	b.a = b.Block.AddValue(Gt, 0, types.Bool, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Le() {
	b.a = b.Block.AddValue(Le, 0, types.Bool, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Ge() {
	b.a = b.Block.AddValue(Ge, 0, types.Bool, b.b, b.a).AddReg(ir.R0)
}

func (b *Builder) Call(fnname string) {
	fn := b.Program.FuncNamed(fnname)
	rettype := b.Program.Types().Func(fn.Sig).ReturnType()
	b.a = b.Block.AddValueAny(Call, 0, rettype, fn).AddReg(ir.R0)
}

func (b *Builder) Jump(label string, id int) {
	b.jump(Jump, label, id)
}

func (b *Builder) jump(op Op, label string, id int, ops ...ir.Value) {
	name := label + strconv.Itoa(id)
	if bid, found := b.labels[name]; found {
		b.Block.UpdateTerminator(op, append(ops, bid)...)
		return
	}

	b.Block.UpdateTerminator(op, append(ops, ir.Value{})...)

	if b.refs == nil {
		b.refs = make(map[string][]ref)
	}
	b.refs[name] = append(b.refs[name], ref{b.Block.Terminator(), len(ops)})
}

func (b *Builder) JumpIfFalse(label string, id int) {
	b.jump(If, label, id, b.a, ir.Value{})
}

func (b *Builder) JumpToEpilogue() {
	b.jump(Jump, b.Func.Name+".epilogue", 0)
}

func (b *Builder) Label(label string, id int) {
	name := label + strconv.Itoa(id)
	bid := b.Func.NewBlock(name, Jump, 0, ir.Value{})

	if b.Block != nil {
		term := b.Block.Terminator()
		if term.Op() == If {
			ops := term.Operands()
			ops[1] = bid
			term.SetOperands(ops...)
		} else {
			lastOper := term.NumOperands() - 1
			if term.Operand(lastOper).IsNil() {
				term.SetOperand(lastOper, bid)
			}
		}
	}

	b.Block = bid.Block()

	// fixup any references to this label
	if refs, found := b.refs[name]; found {
		for _, ref := range refs {
			ref.instr.SetOperand(ref.oper, bid)
		}
		delete(b.refs, name)
	}

	if b.labels == nil {
		b.labels = make(map[string]ir.Value)
	}

	// record the label for future references
	b.labels[name] = bid
}
