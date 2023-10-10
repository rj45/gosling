package ir

import (
	"fmt"
	"io"
	"strings"

	"github.com/rj45/gosling/types"
)

// ValueID is a reference to a value in the IR. Operands of values are
// referenced by ValueID, and any references to blocks are also ValueIDs
// pointing to the terminator instruction of the block.
type ValueID uint32

// InvalidValue indicates an invalid or missing value.
const InvalidValue ValueID = 0

// value represents the most frequently required data about an AST node
// or IR value, stored bit-packed into a 64 bit integer.
//
// It is 64 bits long:
// 10 bits for the OpID
// 12 bits for the TypeID
// 14 bits for the blockID
// 8 bits for the number of operands
// 20 bits for the index of the first operand
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

func (v ValueID) string(f *Func) string {
	if b, ok := f.blockForValue[v]; ok {
		return f.block[b].Name
	}
	if f.Regs(v).IsEmpty() {
		return fmt.Sprintf("v%d", int(v))
	}
	return f.Regs(v).String()
}

func (v ValueID) dump(w io.Writer, f *Func) {
	var oper []string

	for i := 0; i < f.NumOperands(v); i++ {
		ov := f.Operand(v, i)
		if c, ok := f.valueConstant[ov]; ok {
			oper = append(oper, c.String())
			continue
		}

		if !f.Regs(ov).IsEmpty() {
			oper = append(oper, f.Regs(ov).String())
			continue
		}

		oper = append(oper, ov.string(f))
	}

	ostr := ""
	if len(oper) > 0 {
		ostr = " " + strings.Join(oper, ", ")
	}

	val := f.value[v]
	if f.typ[val.typeID()] == types.Void {
		fmt.Fprintf(w, "\t%s%s\n", val.op(), ostr)
		return
	}

	fmt.Fprintf(w, "\t%s = %s%s\n", v.string(f), f.Op(v), ostr)
}
