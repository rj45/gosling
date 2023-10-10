package ir

import (
	"fmt"
	"io"
	"strings"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// blockID is the ID of a block. Usually blocks are identified
// with a ValueID, but internally blockID is used for the block's
// index.
type blockID uint32

// Block is a basic block in the function. Blocks are not in any
// specific order until the blocks are scheduled, at which point
// the blocks are in order. The successors of blocks can be
// determined by looking at the operands of the terminator.
//
// Blocks have a list of parameters, which are the Phi instructions
// when the IR is in SSA form. These act as parallel copies, which
// means that all the phi instructions are executed at the same time
// just before the block begins execution, and so their order does
// not matter.
type Block struct {
	*Func

	// Name or label of the block.
	Name string

	// The list of values in the block. These are not in any
	// specific order until the block is scheduled. This list
	// may include InvalidValue values where values have been
	// deleted. New values are inserted in the first InvalidValue
	// slot or are added to the end of the list.
	value []ValueID

	// The terminating instruction of the block.
	// This is generally a branch or return instruction.
	terminator ValueID
}

// NewBlock creates a new block in the function.
func (fn *Func) NewBlock(name string, op Op, token token.Token, args ...ValueID) ValueID {
	id := blockID(len(fn.block))
	if fn.blockForValue == nil {
		fn.blockForValue = make(map[ValueID]blockID)
	}
	terminator := fn.addValue(op, id, fn.lookupType(types.Void), token, args...)
	fn.block = append(fn.block, Block{
		Func:       fn,
		Name:       name,
		terminator: terminator,
	})
	block := &fn.block[id]
	fn.blockForValue[terminator] = id

	return block.ValueID()
}

// NumBlocks returns the number of blocks in the function.
func (fn *Func) NumBlocks() int {
	return len(fn.block)
}

// Block returns the block with the given ID.
func (fn *Func) Block(id ValueID) *Block {
	return &fn.block[fn.value[id].block()]
}

// BlockAt returns the block at the given index.
func (fn *Func) BlockAt(index int) *Block {
	return &fn.block[index]
}

// NumValues returns the number of values in the block.
func (b *Block) NumValues() int {
	return len(b.value)
}

// ValueAt returns the value at the given index.
// Note: this can return an InvalidValue if the value
// has been deleted.
func (b *Block) ValueAt(index int) ValueID {
	return b.value[index]
}

// AddValue adds a new value to the block.
func (b *Block) AddValue(op Op, token token.Token, typ types.Type, args ...ValueID) ValueID {
	// todo: optimize this... in some cases it may be faster to search in b.Func.block
	blockID := b.Func.blockForValue[b.ValueID()]

	id := b.Func.addValue(op, blockID, b.Func.lookupType(typ), token, args...)
	b.value = append(b.value, id)
	return id
}

// AppendValue appends a new value to the block.
func (b *Block) AppendValue(id ValueID) {
	b.value = append(b.value, id)
}

// ValueID returns the ValueID of the block. This is
// currently the same as the terminator value, but this
// may change in the future, so use this function instead.
func (b *Block) ValueID() ValueID {
	return b.terminator
}

// Terminator returns the terminator of the block.
func (b *Block) Terminator() ValueID {
	return b.terminator
}

func (b *Block) UpdateTerminator(Op Op, args ...ValueID) {
	b.Func.SetValueOp(b.terminator, Op)
	b.Func.SetOperands(b.terminator, args...)
}

func (b *Block) Dump() string {
	w := &strings.Builder{}
	b.dump(w)
	return w.String()
}

func (b *Block) dump(w io.Writer) {
	fmt.Fprintln(w, b.Name+":")
	for _, v := range b.value {
		v.dump(w, b.Func)
	}

	if b.terminator != InvalidValue {
		b.terminator.dump(w, b.Func)
	}
}
