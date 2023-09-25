package vm

type Opcode uint8

func (i Instr) Opcode() Opcode {
	return Opcode(i & 0xff)
}

func (i Instr) Arg() int {
	return int((i >> 8) & 0xffffffff)
}

const (
	Undef Opcode = iota
	Prologue
	Load
	Store
	Push
	Pop
	LoadLocal
	LoadInt
	LocalAddr
	Add
	Sub
	Mul
	Div
	Neg
	Eq
	Ne
	Lt
	Le
	Gt
	Ge
	Call
	JumpIfFalse
	Jump
	Return
	Exit
)

var opcodeNames = [...]string{
	Undef:       "undef",
	Prologue:    "prologue",
	Load:        "load",
	Store:       "store",
	Push:        "push",
	Pop:         "pop",
	LoadLocal:   "loadlocal",
	LoadInt:     "loadint",
	LocalAddr:   "localaddr",
	Add:         "add",
	Sub:         "sub",
	Mul:         "mul",
	Div:         "div",
	Neg:         "neg",
	Eq:          "eq",
	Ne:          "ne",
	Lt:          "lt",
	Le:          "le",
	Gt:          "gt",
	Ge:          "ge",
	Call:        "call",
	JumpIfFalse: "jumpiffalse",
	Jump:        "jump",
	Return:      "return",
	Exit:        "exit",
}

func (o Opcode) String() string {
	if o >= Undef && o < Opcode(len(opcodeNames)) {
		return opcodeNames[o]
	}
	return "unknown"
}

var opcodeHasArg = [...]bool{
	Undef:       false,
	Prologue:    true,
	Load:        false,
	Store:       false,
	Push:        false,
	Pop:         false,
	LoadLocal:   true,
	LoadInt:     true,
	LocalAddr:   true,
	Add:         false,
	Sub:         false,
	Mul:         false,
	Div:         false,
	Neg:         false,
	Eq:          false,
	Ne:          false,
	Lt:          false,
	Le:          false,
	Gt:          false,
	Ge:          false,
	Call:        true,
	JumpIfFalse: true,
	Jump:        true,
	Return:      false,
	Exit:        false,
}
