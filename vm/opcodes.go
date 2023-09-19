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
	JumpIfFalse
	Jump
	Epilogue
)
