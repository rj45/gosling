package ir

import (
	"fmt"
	"io"
	"strings"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// Func is a function.
type Func struct {
	// The file the func is in, which is also a reference
	// to the source code of the file, which all tokens
	// reference.
	*token.File

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

	// The list of constants
	constantValue map[Constant]ValueID
	valueConstant map[ValueID]Constant
}

// NewFunc creates a new Func.
func NewFunc(file *token.File) *Func {
	if file == nil {
		panic("file cannot be nil")
	}
	return &Func{
		File: file,
		value: []value{
			// the first value is always an invalid value
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
		block: []Block{
			// the first block is always an invalid block
			{},
		},
	}
}

// NumValues returns the number of values in the function.
func (fn *Func) NumValues() int {
	return len(fn.value) - 1
}

// Value returns the Value for the given id.
func (fn *Func) Value(id ValueID) Value {
	return fn.valueForID(id)
}

// ValueAt returns the Value at the given index.
func (fn *Func) ValueAt(index int) Value {
	return fn.valueForID(ValueID(index + 1))
}

func (fn *Func) valueForID(id ValueID) Value {
	val := fn.value[id]
	val &^= value(0xfffff) << 44
	val = val | value(id)<<44
	v := Value{value: val, fn: fn}
	return v
}

// AddValue adds a new value to the Func without a block.
func (fn *Func) AddValue(op Op, token token.Token, typ types.Type, operands ...Value) Value {
	return fn.addValue(op, 0, fn.lookupType(typ), token, operands...)
}

// AddValueAny adds a new value to the Func without a block, with any type of operand.
func (fn *Func) AddValueAny(op Op, token token.Token, typ types.Type, operands ...any) Value {
	values := make([]Value, len(operands))
	for i, oper := range operands {
		values[i] = fn.valueForAny(oper)
	}
	return fn.addValue(op, 0, fn.lookupType(typ), token, values...)
}

func (fn *Func) addValue(op Op, block BlockID, typ typeID, token token.Token, operands ...Value) Value {
	id := ValueID(len(fn.value))

	firstOperand := len(fn.operand)
	for i, oper := range operands {
		fn.operand = append(fn.operand, oper.id())
		if fn.operand[firstOperand+i] != oper.id() {
			panic("invalid operand")
		}
	}

	fn.value = append(fn.value, newValue(op.OpID(), typ, block, len(operands), firstOperand))
	fn.token = append(fn.token, token)
	fn.regs = append(fn.regs, 0)

	return fn.valueForID(id)
}

func (fn *Func) valueForAny(v any) Value {
	switch v := v.(type) {
	case ValueID:
		return fn.valueForID(v)
	case Value:
		return v
	case Constant:
		return fn.ValueForConst(v)
	case RegID:
		c := RegConst(NewRegMask(v))
		return fn.ValueForConst(c)
	case RegMask:
		c := RegConst(v)
		return fn.ValueForConst(c)
	case int:
		c := IntConst(int64(v))
		return fn.ValueForConst(c)
	case int64:
		c := IntConst(v)
		return fn.ValueForConst(c)
	case *Func:
		c := FuncConst(v)
		return fn.ValueForConst(c)
	case bool:
		c := BoolConst(v)
		return fn.ValueForConst(c)
	}
	panic(fmt.Sprintf("invalid value: %T", v))
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

// String returns the name of the function.
func (fn *Func) String() string {
	return fn.Name
}

// ValueForConst returns the value for the given constant.
func (fn *Func) ValueForConst(c Constant) Value {
	if id, ok := fn.constantValue[c]; ok {
		return fn.valueForID(id)
	}
	if fn.constantValue == nil {
		fn.constantValue = map[Constant]ValueID{}
		fn.valueConstant = map[ValueID]Constant{}
	}
	val := fn.AddValue(Const, 0, types.Int)
	fn.constantValue[c] = val.id()
	fn.valueConstant[val.id()] = c

	return val
}

// Dump returns a string representation of the function.
func (fn *Func) Dump() string {
	w := &strings.Builder{}
	fn.dump(w)
	return w.String()
}

func (fn *Func) dump(w io.Writer) {
	sig := fn.Types().Func(fn.Sig)
	params := make([]string, len(sig.ParamTypes()))
	for i, typ := range sig.ParamTypes() {
		params[i] = fmt.Sprintf("%s %s", RegID(i).String(), fn.Types().StringOf(typ))
	}
	pstr := "(" + strings.Join(params, ", ") + ")"
	if sig.ReturnType() != types.Void {
		pstr += " " + fn.Types().StringOf(sig.ReturnType())
	}
	fmt.Fprintln(w, "func", fn.Name+pstr, "{")
	for _, block := range fn.block {
		if block.id == InvalidBlock {
			continue
		}
		block.dump(w)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
