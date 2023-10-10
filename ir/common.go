package ir

type CommonOp uint8

const (
	Invalid CommonOp = iota
	Const
	Reg
)

var opNames = [...]string{
	Invalid: "Invalid",
	Const:   "Const",
	Reg:     "Reg",
}

func (op CommonOp) String() string {
	if op >= CommonOp(len(opNames)) {
		panic("invalid op")
	}
	name := opNames[op]
	if name == "" {
		panic("missing name for op")
	}
	return name
}

func (op CommonOp) Category() OpCategory {
	if op == Const {
		return ConstOp
	}
	return InvalidOp
}

func (op CommonOp) OpID() OpID {
	return OpID(op)
}

var _ = RegisterOpSet(Common, func(oi OpID) Op {
	return CommonOp(oi)
})
