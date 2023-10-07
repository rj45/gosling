package types

type TypeKind uint8

const (
	BasicType TypeKind = iota
	FuncType
)

// Type identifies a type within the universe of types.
//
// It is 24 bits long:
// 4 bits for the TypeKind
// 2 bits for the number of pointer indirections
// 18 bits for the index into the universe of types
type Type uint32

func newType(kind TypeKind, index int, indirections int) Type {
	if kind < BasicType || kind > FuncType {
		panic("kind out of range")
	}
	if index < 0 || index > 0x3ffff {
		panic("index must be between 0 and 0x3ffff")
	}
	if indirections < 0 || indirections > 3 {
		panic("indirections must be between 0 and 3")
	}

	return Type((uint32(kind) << 20) | (uint32(indirections&3) << 18) | uint32(index))
}

// Kind returns the kind of type.
func (t Type) Kind() TypeKind {
	return TypeKind(t >> 20)
}

// Index into the universe of types.
func (t Type) Index() int {
	return int(t & 0x3ffff)
}

// Indirections returns the number of pointer indirections.
// That is, `p` is 0, `*p` is 1, `**p` is 2, etc.
// The maximum number of indirections is 3.
func (t Type) Indirections() int {
	return int((t >> 18) & 3)
}

func (t Type) Deref() Type {
	if t.Indirections() == 0 {
		panic("cannot deref non-pointer type")
	}
	return newType(t.Kind(), t.Index(), t.Indirections()-1)
}

func (t Type) Pointer() Type {
	if t.Indirections() == 3 {
		panic("cannot take pointer of triple pointer type")
	}
	return newType(t.Kind(), t.Index(), t.Indirections()+1)
}
