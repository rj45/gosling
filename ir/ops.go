package ir

type Opcode uint32

type Op interface {
	Opcode() Opcode
	IsControlFlow() bool
	IsDataFlow() bool
	StartsBlock() bool
	EndsBlock() bool
	String() string
}

const (
	OpInvalid Opcode = iota

	// Control flow nodes
	OpStart
	OpRegion
	OpReturn
	OpIf
	OpThen
	OpElse
	OpStop

	// Data+Control merge node
	OpPhi

	// Data flow nodes
	OpError
	OpConst
	OpCopy
	OpExtract // Extracts a value from a tuple
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpNeg
	// ...
)

var _ Op = Opcode(0)

func (o Opcode) Opcode() Opcode { return o }
func (o Opcode) IsControlFlow() bool {
	switch o {
	case OpStart, OpRegion, OpReturn, OpIf, OpThen, OpElse, OpStop:
		return true
	}
	return false
}
func (o Opcode) IsDataFlow() bool {
	return !o.IsControlFlow()
}

func (o Opcode) StartsBlock() bool {
	switch o {
	case OpStart, OpRegion, OpStop:
		return true
	}
	return false
}

func (o Opcode) EndsBlock() bool {
	switch o {
	case OpReturn, OpIf:
		return true
	}
	return false
}

func (o Opcode) String() string {
	switch o {
	case OpStart:
		return "Start"
	case OpRegion:
		return "Region"
	case OpReturn:
		return "Return"
	case OpIf:
		return "If"
	case OpThen:
		return "Then"
	case OpElse:
		return "Else"
	case OpStop:
		return "Stop"
	case OpPhi:
		return "Phi"
	case OpError:
		return "Error"
	case OpConst:
		return "Const"
	case OpCopy:
		return "Copy"
	case OpExtract:
		return "Extract"
	case OpAdd:
		return "Add"
	case OpSub:
		return "Sub"
	case OpMul:
		return "Mul"
	case OpDiv:
		return "Div"
	case OpNeg:
		return "Neg"
	}
	return "Invalid"
}
