package ir

import (
	"fmt"
	"io"
	"strings"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// ValueID is a reference to a value in the IR. Operands of values are
// referenced by ValueID, and any references to blocks are also ValueIDs
// pointing to the terminator instruction of the block.
type ValueID uint32

// InvalidValue indicates an invalid or missing value.
const InvalidValue ValueID = 0

// Value is an IR value. This is *Func, ValueID pair.
type Value struct {
	value
	fn *Func
}

// ID of the Value.
func (v Value) ID() ValueID {
	return ValueID(v.id())
}

// Func returns the function the value is in.
func (v Value) Func() *Func {
	return v.fn
}

// IsNil returns true if the value is nil/empty/InvalidValue.
func (v Value) IsNil() bool {
	return v.fn == nil || v.id() == InvalidValue
}

// OpID of the Value.
func (v Value) OpID() OpID {
	return v.opID()
}

// Op of the Value.
func (v Value) Op() Op {
	return v.op()
}

// SetOp sets the Op of the Value and returns a new Value.
func (v Value) SetOp(op Op) Value {
	v.fn.value[v.id()] = newValue(op.OpID(), v.typeID(), v.block(), v.numOperands(), v.fn.value[v.id()].firstOperand())
	return v.fn.valueForID(v.id())
}

// Token of the Value.
func (v Value) Token() token.Token {
	return v.fn.token[v.id()]
}

// SetToken sets the Token of the Value and returns a new Value.
func (v Value) SetToken(tok token.Token) Value {
	v.fn.token[v.id()] = tok
	return v
}

// Type of the Value.
func (v Value) Type() types.Type {
	return v.fn.typ[v.typeID()]
}

// SetType sets the Type of the Value and returns a new Value.
func (v Value) SetType(typ types.Type) Value {
	id := v.fn.lookupType(typ)
	v.fn.value[v.id()] = newValue(v.opID(), id, v.block(), v.numOperands(), v.fn.value[v.id()].firstOperand())
	return v.fn.valueForID(v.id())
}

// Block of the Value.
func (v Value) Block() *Block {
	return &v.fn.block[v.block()]
}

// SetBlock sets the Block of the Value and returns a new Value.
func (v Value) SetBlock(block *Block) Value {
	if v.Block() != block && v.Block() != nil {
		v.Block().removeValue(v.id())
	}
	v.fn.value[v.id()] = newValue(v.opID(), v.typeID(), block.blockID(), v.numOperands(), v.fn.value[v.id()].firstOperand())
	block.appendValue(v.id())
	return v.fn.valueForID(v.id())
}

// NumOperands returns the number of operands.
func (v Value) NumOperands() int {
	return v.numOperands()
}

// OperandID returns the ValueID of the operand at the given index.
func (v Value) OperandID(index int) ValueID {
	return v.Operand(index).id()
}

// Operand returns the operand at the given index.
func (v Value) Operand(index int) Value {
	if index < 0 || index >= v.numOperands() {
		panic("invalid index")
	}
	return v.fn.valueForID(v.fn.operand[v.fn.value[v.id()].firstOperand()+index])
}

// Operands returns the operands of the value.
func (v Value) Operands() []Value {
	operands := make([]Value, v.numOperands())
	for i := range operands {
		operands[i] = v.Operand(i)
	}
	return operands
}

// FillOperands fills the operands slice with the operands of the value.
func (v Value) FillOperands(operands []Value) []Value {
	operands = operands[:0]
	for i := range operands {
		operands = append(operands, v.Operand(i))
	}
	return operands
}

// SetOperand sets the operand at the given index.
func (v Value) SetOperand(index int, operand Value) {
	v.setOperandID(index, operand.id())
}

// SetOperandID sets the operand ValueID at the given index.
func (v Value) SetOperandID(index int, operand ValueID) {
	v.setOperandID(index, operand)
}

func (v Value) setOperandID(index int, operand ValueID) {
	n := v.fn.value[v.id()]

	if index < 0 || index >= n.numOperands() {
		panic("invalid index")
	}
	v.fn.operand[n.firstOperand()+index] = operand
}

// SetOperands replaces the operands of the value.
func (v Value) SetOperands(operands ...Value) Value {
	fn := v.fn
	id := v.id()
	n := fn.value[id]

	if len(operands) <= n.numOperands() {
		// zero out excess operands
		for i := len(operands); i < n.numOperands(); i++ {
			fn.operand[n.firstOperand()+i] = InvalidValue
			// todo: add to free list
		}
		for i := range operands {
			fn.operand[n.firstOperand()+i] = operands[i].id()
		}
		fn.value[id] = newValue(n.opID(), n.typeID(), n.block(), len(operands), n.firstOperand())
		return fn.valueForID(id)
	}

	// todo: find a chunk of free operands on the free list

	firstOperand := len(fn.operand)
	for _, oper := range operands {
		fn.operand = append(fn.operand, oper.id())
	}

	// zero out the former operands
	// important: do this after appending to fn.operand or
	// the passed in operands could be zeroed out
	for i := 0; i < n.numOperands(); i++ {
		fn.operand[n.firstOperand()+i] = InvalidValue
		// todo: add to free list
	}

	fn.value[id] = newValue(n.opID(), n.typeID(), n.block(), len(operands), firstOperand)
	return fn.valueForID(id)
}

// SetOperandsAny replaces the operands of the value with any type of operand.
func (v Value) SetAnyOperands(operands ...any) Value {
	values := make([]Value, len(operands))
	for i, oper := range operands {
		values[i] = v.fn.valueForAny(oper)
	}
	return v.SetOperands(values...)
}

// IsConstant returns true if the value is a constant.
func (v Value) IsConstant() bool {
	_, found := v.fn.valueConstant[v.id()]
	return found
}

// Constant returns the constant value.
func (v Value) Constant() Constant {
	c, found := v.fn.valueConstant[v.id()]
	if !found {
		panic("not a constant")
	}
	return c
}

// HasRegister returns true if the value has a register.
func (v Value) HasRegister() bool {
	return !v.Regs().IsEmpty()
}

// Regs returns the registers of the value, checking if the value is a register constant as well.
func (v Value) Regs() RegMask {
	if v.IsConstant() {
		if reg, ok := RegValue(v.Constant()); ok {
			return reg
		}
	}
	return v.fn.regs[v.id()]
}

// SetRegs sets the registers of the value.
func (v Value) SetRegs(regs RegMask) Value {
	v.fn.regs[v.id()] = regs
	return v
}

// AddReg adds a register to the value.
func (v Value) AddReg(reg RegID) Value {
	v.fn.regs[v.id()] |= 1 << reg
	return v
}

// IsBlock returns true if the value is a block.
func (v Value) IsBlock() bool {
	_, found := v.fn.blockForValue[v.id()]
	return found
}

// BlockValue returns the block value, if it's a Block.
func (v Value) BlockValue() *Block {
	b, found := v.fn.blockForValue[v.id()]
	if !found {
		panic("not a block")
	}
	return &v.fn.block[b]
}

// String returns the name of the value.
func (v Value) String() string {
	fn := v.fn
	if b, ok := fn.blockForValue[v.id()]; ok {
		return fn.block[b].Name
	}
	if v.IsConstant() {
		return v.Constant().String()
	}
	if fn.regs[v.id()] != 0 {
		return fn.regs[v.id()].String()
	}
	return fmt.Sprintf("v%d", int(v.id()))
}

func (v Value) Dump() string {
	w := &strings.Builder{}
	v.dump(w)
	return w.String()
}

func (v Value) dump(w io.Writer) {
	var oper []string

	for i := 0; i < v.numOperands(); i++ {
		oper = append(oper, v.Operand(i).String())
	}

	ostr := ""
	if len(oper) > 0 {
		ostr = " " + strings.Join(oper, ", ")
	}

	if v.fn.typ[v.typeID()] == types.Void {
		fmt.Fprintf(w, "\t%s%s\n", v.op(), ostr)
		return
	}

	fmt.Fprintf(w, "\t%s = %s%s\n", v.String(), v.Op(), ostr)
}

// value represents the most frequently required data about an AST node
// or IR value, stored bit-packed into a 64 bit integer.
//
// It is 64 bits long:
// 10 bits for the OpID
// 12 bits for the TypeID
// 14 bits for the blockID
// 8 bits for the number of operands
// 20 bits for the index of the first operand (or id if in the Value struct)
type value uint64

// typeID is the index of the type in the function's type table
type typeID uint16

// newValue creates a new AST node or IR value
func newValue(op OpID, typ typeID, block blockID, numOperands int, firstOperand int) value {
	if uint(op) > 0x3ff {
		panic("op out of range")
	}
	if uint(typ) > 0xfff {
		panic("type out of range")
	}
	if uint(block) > 0x3fff {
		panic("block out of range")
	}
	if uint(numOperands) > 0xff {
		panic("too many operands")
	}
	if uint(firstOperand) > 0xfffff {
		panic("firstOperand out of range")
	}

	return value(op) | value(typ)<<10 | value(block)<<22 | value(numOperands)<<36 | value(firstOperand)<<44
}

// opID of the Value
func (n value) opID() OpID {
	// least significant 10 bits
	return OpID(n & 0x3ff)
}

// op of the value
func (n value) op() Op {
	return opSets[(n>>8)&3][n&0xff]
}

// typeID of the Value
func (n value) typeID() typeID {
	return typeID((n >> 10) & 0xfff)
}

// block of the Value
func (n value) block() blockID {
	return blockID(int(n>>22) & 0x3fff)
}

// numOperands returns the number of operands
func (n value) numOperands() int {
	return int(n>>36) & 0xff
}

// firstOperand returns the first operand in the operands list
func (n value) firstOperand() int {
	return int(n>>44) & 0xfffff
}

// id returns the ValueID of the value (if in the Value struct)
func (n value) id() ValueID {
	return ValueID(int(n>>44) & 0xfffff)
}
