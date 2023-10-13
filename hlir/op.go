package hlir

import "github.com/rj45/gosling/ir"

type Op uint16

const (
	Invalid Op = iota + Op(ir.HLIR)<<8

	Prologue
	Epilogue

	// Stack operators
	Push
	Pop
	LoadLocal
	StoreLocal
	Load
	Store
	LocalAddr

	// Constant operators
	LoadInt

	// Binary operators
	Add
	Sub
	Mul
	Div

	// Unary operators
	Neg
	Move

	// Comparison operators
	Eq
	Ne
	Lt
	Gt
	Le
	Ge

	// Other operators
	Addr
	Deref
	Call

	// Control flow operators
	Jump
	If
	Return
)

var opNames = [...]string{
	Invalid:    "Invalid",
	Prologue:   "Prologue",
	Epilogue:   "Epilogue",
	Push:       "Push",
	Pop:        "Pop",
	LoadLocal:  "LoadLocal",
	StoreLocal: "StoreLocal",
	Load:       "Load",
	Store:      "Store",
	LocalAddr:  "LocalAddr",
	LoadInt:    "LoadInt",
	Add:        "Add",
	Sub:        "Sub",
	Mul:        "Mul",
	Div:        "Div",
	Neg:        "Neg",
	Move:       "Move",
	Eq:         "Eq",
	Ne:         "Ne",
	Lt:         "Lt",
	Gt:         "Gt",
	Le:         "Le",
	Ge:         "Ge",
	Addr:       "Addr",
	Deref:      "Deref",
	Call:       "Call",
	Jump:       "Jump",
	If:         "If",
	Return:     "Return",
}

func (op Op) String() string {
	if op >= Op(len(opNames)) {
		return "Invalid"
	}
	name := opNames[op]
	if name == "" {
		panic("missing name for op")
	}
	return name
}

func (op Op) Category() ir.OpCategory {
	return 0
}

func (op Op) OpID() ir.OpID {
	return ir.OpID(op)
}

var _ = ir.RegisterOpSet(ir.HLIR, func(oi ir.OpID) ir.Op {
	return Op(oi)
})
