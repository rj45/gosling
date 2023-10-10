package ir

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// Func is a function.
type Func struct {
	// The file the func is in, which is also a reference
	// to the source code of the file, which all tokens
	// reference.
	*File

	// The program the function is in.
	*Program

	// The Name of the function
	Name string

	// Sig is the type signature of the function.
	Sig types.Type

	// The list of values in the block, indexed by ValueID
	value []value

	// The list of types indexed by ValueID.
	typ []types.Type

	// The list of assigned registers indexed by ValueID.
	regs []RegMask

	// The list of operands for values.
	operand []ValueID

	// The list of tokens indexed by ValueID. Tokens also represent
	// the position of the value in the source code, as well as any
	// name or literal associated with the value.
	token []token.Token

	// The list of basic blocks in the function.
	block []Block

	// The list of blocks indexed by ValueID.
	blockForValue map[ValueID]blockID

	// The list of constants
	constantValue map[Constant]ValueID
	valueConstant map[ValueID]Constant
}

// NewFunc creates a new Func.
func NewFunc(file *File) *Func {
	if file == nil {
		panic("file cannot be nil")
	}
	return &Func{
		File: file,
		value: []value{
			// the first value is always an illegal value
			newValue(0, 0, 0, 0, 0),
		},
		token: []token.Token{
			// the first token is always an illegal token
			token.Token(0),
		},
		typ: []types.Type{
			// the first types are always None and Void
			types.None, types.Void,
		},
		regs: []RegMask{0},
	}
}

// Token returns the token for the given node.
func (fn *Func) Token(id ValueID) token.Token {
	return fn.token[id]
}

// OpID returns the ID of the Op for the given value.
func (fn *Func) OpID(id ValueID) OpID {
	return fn.value[id].opID()
}

// Returns the Op for the given value.
// This has been optimized to not allocate a new Op.
func (fn *Func) Op(id ValueID) Op {
	return fn.value[id].op()
}

// BlockFor returns the block of the given value.
func (fn *Func) BlockFor(id ValueID) ValueID {
	return fn.block[fn.value[id].block()].ValueID()
}

// NumValues returns the number of values in the function.
func (fn *Func) NumValues() int {
	return len(fn.value)
}

// Type returns the type of the given node
func (fn *Func) Type(id ValueID) types.Type {
	if id >= ValueID(len(fn.typ)) {
		return types.None
	}
	return fn.typ[id]
}

// SetType sets the type of the given node
func (fn *Func) SetType(id ValueID, typ types.Type) {
	for id >= ValueID(len(fn.typ)) {
		fn.typ = append(fn.typ, types.None)
	}
	fn.typ[id] = typ
}

// NumOperands returns the number of operands for the given value
func (fn *Func) NumOperands(id ValueID) int {
	return fn.value[id].numOperands()
}

// Operand returns the nth operand of the given node
func (fn *Func) Operand(id ValueID, index int) ValueID {
	n := fn.value[id]
	if index < 0 || index >= n.numOperands() {
		panic("invalid index")
	}
	return fn.operand[n.firstOperand()+index]
}

// Operands returns the operands of a node
func (fn *Func) Operands(id ValueID) []ValueID {
	n := fn.value[id]

	start := n.firstOperand()
	end := n.numOperands() + start
	return fn.operand[start:end]
}

// SetOperand sets the nth operand of the given node
func (fn *Func) SetOperand(id ValueID, index int, operand ValueID) {
	n := fn.value[id]

	if index < 0 || index >= n.numOperands() {
		panic("invalid index")
	}
	fn.operand[n.firstOperand()+index] = operand
}

// SetOperands sets the operands of a node.
func (fn *Func) SetOperands(id ValueID, operands ...ValueID) {
	n := fn.value[id]

	if len(operands) <= n.numOperands() {
		// zero out excess operands
		for i := len(operands); i < n.numOperands(); i++ {
			fn.operand[n.firstOperand()+i] = InvalidValue
			// todo: add to free list
		}

		copy(fn.operand[n.firstOperand():], operands)
		fn.value[id] = newValue(n.opID(), n.typeID(), n.block(), len(operands), n.firstOperand())
		return
	}

	// todo: find a chunk of free operands on the free list

	firstOperand := len(fn.operand)
	fn.operand = append(fn.operand, operands...)

	// zero out the former operands
	// important: do this after appending to fn.operand or
	// the passed in operands could be zeroed out
	for i := 0; i < n.numOperands(); i++ {
		fn.operand[n.firstOperand()+i] = InvalidValue
		// todo: add to free list
	}

	fn.value[id] = newValue(n.opID(), n.typeID(), n.block(), len(operands), firstOperand)
}

// AddValue adds a new value to the Func without a block.
func (fn *Func) AddValue(op Op, token token.Token, typ types.Type, operands ...ValueID) ValueID {
	return fn.addValue(op, 0, fn.lookupType(typ), token, operands...)
}

func (fn *Func) addValue(op Op, block blockID, typ typeID, token token.Token, operands ...ValueID) ValueID {
	id := ValueID(len(fn.value))

	firstOperand := len(fn.operand)
	fn.operand = append(fn.operand, operands...)

	fn.value = append(fn.value, newValue(op.OpID(), typ, block, len(operands), firstOperand))
	fn.token = append(fn.token, token)
	fn.regs = append(fn.regs, 0)

	return id
}

func (fn *Func) lookupType(typ types.Type) typeID {
	if typ == types.None {
		return 0
	}

	// todo: use a map?
	for i, t := range fn.typ {
		if t == typ {
			return typeID(i)
		}
	}

	id := typeID(len(fn.typ))
	fn.typ = append(fn.typ, typ)
	return id
}

// SetValueBlock sets the block of the given value.
func (fn *Func) SetValueBlock(id ValueID, block ValueID) {
	n := fn.value[id]
	fn.value[id] = newValue(n.opID(), n.typeID(), fn.blockForValue[block], n.numOperands(), n.firstOperand())
}

// SetValueOp sets the op of the given value.
func (fn *Func) SetValueOp(id ValueID, op Op) {
	n := fn.value[id]
	fn.value[id] = newValue(op.OpID(), n.typeID(), n.block(), n.numOperands(), n.firstOperand())
}

// String returns the name of the function.
func (fn *Func) String() string {
	return fn.Name
}

func (fn *Func) ValueBytes(id ValueID) []byte {
	return fn.TokenBytes(fn.Token(id))
}

func (fn *Func) ValueString(id ValueID) string {
	return fn.TokenString(fn.Token(id))
}

// ValueForConst returns the value for the given constant.
func (fn *Func) ValueForConst(c Constant) ValueID {
	if id, ok := fn.constantValue[c]; ok {
		return id
	}
	if fn.constantValue == nil {
		fn.constantValue = map[Constant]ValueID{}
		fn.valueConstant = map[ValueID]Constant{}
	}
	id := fn.AddValue(Const, 0, types.Int)
	fn.constantValue[c] = id
	fn.valueConstant[id] = c

	return id
}

// Regs returns the registers for the given value.
func (fn *Func) Regs(id ValueID) RegMask {
	return fn.regs[id]
}

// SetRegs sets the registers for the given value.
func (fn *Func) SetRegs(id ValueID, regs RegMask) {
	fn.regs[id] = regs
}

// AddReg adds a register to the given value.
func (fn *Func) AddReg(id ValueID, reg RegID) {
	if len(fn.regs) <= int(id) {
		log.Panicf("expected %d to be less than %d", id, len(fn.regs))
	}
	fn.regs[id] |= 1 << reg
}

func (fn *Func) Dump() string {
	w := &strings.Builder{}
	fn.dump(w)
	return w.String()
}

func (fn *Func) dump(w io.Writer) {
	// todo: emit signature
	fmt.Fprintln(w, "func", fn.Name, "{")
	for _, block := range fn.block {
		block.dump(w)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
