package ir

// OpID is the ID of an operation in the IR.
//
// It is 10 bits long:
// 2 bits for the Level
// 8 bits for the Op index at that level
type OpID uint16

// Level returns the abstraction level of the operation.
func (id OpID) Level() OpLevel {
	return OpLevel(id >> 8)
}

// Index returns the index of the operation at the given level.
func (id OpID) Index() uint8 {
	return uint8(id)
}

// NewOpID creates a new OpID.
func NewOpID(level OpLevel, index uint8) OpID {
	if level > 0x3 {
		panic("level out of range")
	}
	if index > 0xff {
		panic("index out of range")
	}
	return OpID(level)<<8 | OpID(index)
}

// String returns the string representation of the operation.
func (id OpID) String() string {
	return opSets[id.Level()][id.Index()].String()
}

// OpLevel is the level of abstraction of an operation.
type OpLevel uint8

const (
	// Common ops -- common ops that are used by all levels.
	Common OpLevel = iota

	// High Level IR -- the result of semantic analysis.
	HLIR

	// Low Level IR -- the result of lowering.
	LLIR

	// Assembly -- the result of code generation.
	ASM
)

// OpCategory is the the general abstract category of an operation.
type OpCategory uint8

const (
	// InvalidOp is an invalid operation.
	InvalidOp OpCategory = iota

	// ConstOp is a constant operation.
	ConstOp
)

// Op is an operation in the IR. This can be an instruction,
// or a node kind in the case of the AST.
type Op interface {
	// OpID returns the ID of the operation.
	OpID() OpID

	// String returns the string representation of the operation.
	String() string

	// Category returns the kind of operation.
	Category() OpCategory

	// ... todo: add more methods for dealing with ops in an
	// OpSet agnostic way.
}

var opSets = [4][256]Op{}

// RegisterOpSet registers a set of operations at the given level.
func RegisterOpSet(level OpLevel, f func(OpID) Op) int {
	for i := 0; i < 256; i++ {
		opSets[level][i] = f(OpID(int(level)*256 + i))
	}
	return 0
}

// OpFromID returns the Op for the given ID.
func OpFromID(id OpID) Op {
	return opSets[id.Level()][id.Index()]
}
