package ir

import (
	"fmt"
	"io"
	"strings"

	"github.com/rj45/gosling/token"
	"github.com/rj45/gosling/types"
)

// BlockID is the ID of a block.
type BlockID uint32

// InvalidBlock is an invalid block ID.
const InvalidBlock BlockID = 0

// Block is a basic block in the function. Blocks are not in any
// specific order until the blocks are scheduled, at which point
// the blocks are in order.
//
// Blocks have a list of parameters, which are the Phi instructions
// when the IR is in SSA form. These act as parallel copies, which
// means that all the phi instructions are executed at the same time
// just before the block begins execution, and so their order does
// not matter.
type Block struct {
	*Func

	id BlockID

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

	// The list of predecessors of the block.
	preds []BlockID

	// The list of successors of the block.
	succs []BlockID
}

// NewBlock creates a new block in the function.
// Warning: this will invalidate previous *Block pointers. Writes to previous
// pointers should redirect to the new block, but reads may return stale data.
func (fn *Func) NewBlock(name string, op Op, token token.Token, args ...Value) *Block {
	id := BlockID(len(fn.block))
	terminator := fn.addValue(op, id, fn.lookupType(types.Void), token, args...)
	fn.block = append(fn.block, Block{
		id:         id,
		Func:       fn,
		Name:       name,
		terminator: terminator.id(),
	})
	block := &fn.block[id]

	return block
}

// NumBlocks returns the number of blocks in the function.
func (fn *Func) NumBlocks() int {
	return len(fn.block) - 1
}

// Block returns the block with the given ID.
func (fn *Func) Block(id BlockID) *Block {
	return &fn.block[id]
}

// BlockAt returns the block at the given index.
func (fn *Func) BlockAt(index int) *Block {
	return &fn.block[index+1]
}

// NumValues returns the number of values in the block.
func (b *Block) NumValues() int {
	return len(b.value)
}

// ValueAt returns the value at the given index.
// Note: this can return an InvalidValue if the value
// has been deleted.
func (b *Block) ValueAt(index int) Value {
	return b.Func.valueForID(b.value[index])
}

// ID returns the ID of the block.
func (b *Block) ID() BlockID {
	return b.id
}

// AddValue adds a new value to the block.
func (b *Block) AddValue(op Op, token token.Token, typ types.Type, args ...Value) Value {
	blockID := b.ID()

	val := b.Func.addValue(op, blockID, b.Func.lookupType(typ), token, args...)
	b.Func.block[blockID].value = append(b.value, val.id())
	return val
}

// AddValueAny adds a new value to the Block, with any type of operand.
func (b *Block) AddValueAny(op Op, token token.Token, typ types.Type, operands ...any) Value {
	blockID := b.ID()

	values := make([]Value, len(operands))
	for i, oper := range operands {
		values[i] = b.Func.valueForAny(oper)
	}

	val := b.Func.addValue(op, blockID, b.Func.lookupType(typ), token, values...)
	b.appendValue(val.id())

	return val
}

func (b *Block) appendValue(id ValueID) {
	b = &b.Func.block[b.ID()] // fix invalid *Block pointers
	b.value = append(b.value, id)
}

// removeValue removes the value from the block.
func (b *Block) removeValue(id ValueID) {
	b = &b.Func.block[b.ID()] // fix invalid *Block pointers
	for i, v := range b.value {
		if v == id {
			b.value[i] = InvalidValue
			return
		}
	}
}

// Terminator returns the terminator of the block.
func (b *Block) Terminator() Value {
	return b.Func.valueForID(b.terminator)
}

// UpdateTerminator updates the op & operands of the terminator.
func (b *Block) UpdateTerminator(Op Op, operands ...Value) {
	b.Func.valueForID(b.terminator).SetOp(Op).SetOperands(operands...)
}

// NumSuccessors returns the number of successors of the block.
func (b *Block) NumSuccessors() int {
	return len(b.succs)
}

// Successor returns the successor at the given index.
func (b *Block) Successor(index int) *Block {
	return b.Func.Block(b.succs[index])
}

// AddSuccessor adds a successor to the block.
// A predecessor edge is also added to the successor block.
func (b *Block) AddSuccessor(succ *Block) {
	b = &b.Func.block[b.ID()] // fix invalid *Block pointers
	b.succs = append(b.succs, succ.ID())
	b.Func.block[succ.ID()].addPredecessor(b)
}

// InsertSuccessor inserts a successor at the given index.
// A predecessor edge is also added to the successor block.
func (b *Block) InsertSuccessor(index int, succ *Block) {
	b = &b.Func.block[b.ID()] // fix invalid *Block pointers
	if index == len(b.succs) {
		b.AddSuccessor(succ)
		return
	}
	if index > len(b.succs) {
		panic("invalid index")
	}
	b.succs = append(b.succs, 0)
	copy(b.succs[index+1:], b.succs[index:])
	b.succs[index] = succ.ID()
	b.Func.block[succ.ID()].addPredecessor(b)
}

// NumPredecessors returns the number of predecessors of the block.
func (b *Block) NumPredecessors() int {
	return len(b.preds)
}

// Predecessor returns the predecessor at the given index.
func (b *Block) Predecessor(index int) *Block {
	return b.Func.Block(b.preds[index])
}

func (b *Block) addPredecessor(pred *Block) {
	b = &b.Func.block[b.ID()] // fix invalid *Block pointers
	b.preds = append(b.preds, pred.ID())
}

// Dump emits the IR for the block as a string.
func (b *Block) Dump() string {
	w := &strings.Builder{}
	b.dump(w)
	return w.String()
}

func (b *Block) dump(w io.Writer) {
	fmt.Fprintln(w, b.Name+":")
	for _, v := range b.value {
		if v == InvalidValue {
			continue
		}
		b.Func.valueForID(v).dump(w)
	}

	if b.terminator != InvalidValue {
		b.Func.valueForID(b.terminator).dump(w)
	}
}
