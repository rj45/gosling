package types

type Type interface {
	Underlying() Type
	String() string
}
