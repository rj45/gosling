package types

type Func struct {
	// todo parameters
	ret Type
}

func (f *Func) String() string {
	if f.ret != nil {
		return "func() " + f.ret.String()
	}
	return "func()"
}

func (f *Func) Underlying() Type {
	return f
}

// ReturnType returns the return type of the function.
func (f *Func) ReturnType() Type {
	return f.ret
}
