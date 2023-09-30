package types

const (
	None Type = iota
	Void
	Int
	Bool
	UntypedInt
)

type basicFlags uint32

const (
	isUntyped basicFlags = 1 << iota
	isInteger
	isBoolean
)

type Basic struct {
	name  string
	flags basicFlags
}

var basicInfos = [...]Basic{
	None:       {"none", 0},
	Void:       {"void", 0},
	Int:        {"int", isInteger},
	Bool:       {"bool", isBoolean},
	UntypedInt: {"int constant", isUntyped | isInteger},
}

func (b Basic) String() string {
	return b.name
}

func (b Basic) IsInteger() bool {
	return b.flags&isInteger != 0
}

func (b Basic) IsBoolean() bool {
	return b.flags&isBoolean != 0
}

func (b Basic) IsUntyped() bool {
	return b.flags&isUntyped != 0
}
