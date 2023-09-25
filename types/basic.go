package types

type Basic uint8

const (
	Invalid Basic = iota
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

type basicInfo struct {
	name  string
	flags basicFlags
}

var basicInfos = [...]basicInfo{
	Invalid:    {"invalid", 0},
	Void:       {"void", 0},
	Int:        {"int", isInteger},
	Bool:       {"bool", isBoolean},
	UntypedInt: {"untyped int", isUntyped | isInteger},
}

func (b Basic) String() string {
	return basicInfos[b].name
}

func (b Basic) IsInteger() bool {
	return basicInfos[b].flags&isInteger != 0
}

func (b Basic) IsBoolean() bool {
	return basicInfos[b].flags&isBoolean != 0
}

func (b Basic) IsUntyped() bool {
	return basicInfos[b].flags&isUntyped != 0
}

func (b Basic) Underlying() Type {
	return b
}
