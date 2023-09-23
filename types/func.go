package types

type Func struct {
	// todo parameters & return value
}

func (f *Func) String() string {
	return "func()"
}

func (f *Func) Underlying() Type {
	return f
}
